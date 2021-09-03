package main

import (
	"context"
	"errors"
	"fmt"
	"gin-web/artisan"
	"gin-web/common"
	"gin-web/component/logs"
	"gin-web/conf"
	"gin-web/crontab"
	"gin-web/kernel"
	"gin-web/queue"
	"gin-web/routes"
	"github.com/PeterYangs/gcmd2"
	"github.com/PeterYangs/tools"
	"github.com/PeterYangs/tools/file/read"
	"github.com/PeterYangs/tools/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"log"
	http_ "net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {

	args := os.Args
	daemon := false
	for k, v := range args {
		if v == "-d" {
			daemon = true
			args[k] = ""
		}
	}

	//直接运行则为阻塞模式，用于开发模式
	if len(args) == 1 {

		args = append(args, "block")

		block(args...)

		return
	}

	switch args[1] {

	case "start":

		//后台运行模式
		if daemon {

			args[1] = "block"
			daemonize(args...)
			return
		}

		args[1] = "block"
		block(args...)

	case "stop":

		err := stop()

		if err != nil {

			log.Println(err)

		}

	case "restart":

		err := stop()

		if err != nil {

			log.Println(err)

			return

		}

		cmd := exec.Command(args[0], "block")
		cmd.Env = os.Environ()
		err = cmd.Start()

		if err != nil {

			log.Println(err)

		}

	case "block":

		serverStart()

	case "artisan":

		artisan.Artisan()

	}

}

func serverStart() {

	kernel.IdInit()

	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	httpOk := make(chan bool)

	httpFail := make(chan bool)

	srv := &http_.Server{}

	go func() {

		sig := <-sigs

		fmt.Println()
		fmt.Println(sig)

		//删除pid文件
		os.Remove("logs/run.pid")

		//删除运行命令
		//os.Remove("logs/cmd")

		c, e := context.WithTimeout(context.Background(), 3*time.Second)

		defer e()

		//关闭http服务
		err := srv.Shutdown(c)

		if err != nil {

			log.Println(err)
		}

		//退出测试
		//f, err := os.OpenFile("logs/quit.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		//
		//f.Write([]byte("http服务已结束," + tools.Date("Y-m-d", time.Now().Unix()) + "\n"))
		//
		//f.Close()

		cancel()
	}()

	go boot(cxt, &wait, httpOk, httpFail)

	//启动http服务
	go httpStart(httpFail, srv)

	<-httpOk

	//等待其他服务退出
	wait.Wait()

	fmt.Println("finish")
}

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}

func logInit(cxt context.Context, wait *sync.WaitGroup) {

	//日志退出标记
	wait.Add(1)

	l := logs.CreateLogs()

	//日志写入任务
	go l.Task(cxt, wait)

}

func httpStart(httpFail chan bool, srv *http_.Server) {

	r := gin.Default()

	//加载配置
	conf.Load()

	//加载全局中间件
	kernel.Load()

	//加载路由
	routes.Load(r)

	//设置端口
	port := os.Getenv("PORT")

	if port == "" {

		port = "8887"
	}

	srv.Addr = ":" + port

	srv.Handler = r

	if err := srv.ListenAndServe(); err != nil && err != http_.ErrServerClosed {

		log.Println(err)

		httpFail <- true

	}

}

func queueInit(cxt context.Context, wait *sync.WaitGroup) {

	//延迟队列的标记
	wait.Add(1)

	for i := 0; i < cast.ToInt(os.Getenv("QUEUE_NUM")); i++ {

		wait.Add(1)

		//启动消息队列
		go queue.Run(cxt, wait)

	}

}

func boot(cxt context.Context, wait *sync.WaitGroup, httpOk chan bool, httpFail chan bool) {

	defer func() {

		httpOk <- true

	}()

	client := http.Client().SetTimeout(1 * time.Second)

	for {

		select {

		case <-httpFail:

			fmt.Println("退出")

			return

		default:

			time.Sleep(200 * time.Millisecond)

			str, err := client.Request().GetToString("http://127.0.0.1:" + os.Getenv("PORT") + "/ping/" + kernel.Id)

			if err == nil && str == "success" {

				//开启任务调度
				go crontab.Run(wait)

				//队列启动
				queueInit(cxt, wait)

				//日志模块初始化
				logInit(cxt, wait)

				//记录pid和启动命令
				runInit()

				return

			}
		}

	}

}

//后台运行
func daemonize(args ...string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	err := cmd.Start()

	if err != nil {

		fmt.Println(err)
	}
}

//阻塞运行
func block(args ...string) {

	sysType := runtime.GOOS

	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if sysType == `linux` {

		runUser := os.Getenv("RUN_USER")

		if runUser == "" || runUser == "nobody" {

			normal(args...)

			return

		}

		cmd := gcmd2.NewCommand("sudo -u "+runUser+" "+tools.Join(" ", args), context.TODO())

		err := cmd.Start()

		if err != nil {

			log.Println(err)

			return

		}

	}

	if sysType == `windows` {

		normal(args...)

		return
	}

}

func normal(args ...string) {

	cmd := gcmd2.NewCommand(tools.Join(" ", args), context.TODO())

	err := cmd.Start()

	if err != nil {

		log.Println(err)
	}

}

func stop() error {

	b, err := common.PathExists("logs/run.pid")

	if err != nil {

		//log.Println(err)

		return err
	}

	if b {

		pid, err := read.Open("logs/run.pid").Read()

		if err != nil {

			//log.Println(err)

			return err

		}

		sysType := runtime.GOOS

		var cmd *exec.Cmd

		if sysType == `windows` {

			cmd = exec.Command("cmd", "/c", "taskkill /f /pid "+string(pid))

		}

		if sysType == `linux` {

			cmd = exec.Command("bash", "-c", "kill "+string(pid))
		}

		err = cmd.Start()

		cmd.Wait()

		//cmd.
		if err != nil {

			//log.Println(err)

			return err

		}

	} else {

		return errors.New("run.pid文件不存在")
	}

	return nil
}

//记录pid和启动命令
func runInit() {

	f, err := os.OpenFile("logs/run.pid", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)

	if err != nil {

		panic(err)
	}

	//记录pid
	f.Write([]byte(cast.ToString(os.Getpid())))

	f.Close()

}
