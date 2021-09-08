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
	"gin-web/redis"
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
	"strings"
	"sync"
	"syscall"
	"time"
)

var isRun = false

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

		args[1] = "block"
		daemonize(args...)

		fmt.Println("starting")

		return

	case "block":

		serverStart()

	case "artisan":

		artisan.Artisan()

	default:

		fmt.Println("命令不存在")

	}

}

//主服务函数
func serverStart() {

	//生成服务id
	kernel.IdInit()

	//检测退出信号
	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//服务退出上下文，主要作用是让其他子组件协程安全退出
	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	httpOk := make(chan bool)

	httpFail := make(chan bool)

	srv := &http_.Server{}

	go func() {

		sig := <-sigs

		fmt.Println()
		fmt.Println(sig)

		if isRun {

			//删除pid文件
			_ = os.Remove("logs/run.pid")

		}

		c, e := context.WithTimeout(context.Background(), 3*time.Second)

		defer e()

		//关闭http服务
		err := srv.Shutdown(c)

		if err != nil {

			log.Println(err)
		}

		//通知子组件协程退出
		cancel()
	}()

	//日志模块初始化
	logInit(cxt, &wait)

	//启动子组件服务
	go boot(cxt, &wait, httpOk, httpFail, sigs)

	//启动http服务
	go httpStart(httpFail, srv)

	//等待http服务启动完成（不论成功或者失败）
	<-httpOk

	//等待其他子组件服务退出
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

//启动日志服务
func logInit(cxt context.Context, wait *sync.WaitGroup) {

	//日志退出标记
	wait.Add(1)

	l := logs.CreateLogs()

	//日志写入任务
	go l.Task(cxt, wait)

}

//启动http服务
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

		//http服务启动失败
		httpFail <- true

	}

}

//启动消息队列服务
func queueInit(cxt context.Context, wait *sync.WaitGroup) {

	//延迟队列的标记
	wait.Add(1)

	for i := 0; i < cast.ToInt(os.Getenv("QUEUE_NUM")); i++ {

		wait.Add(1)

		//启动消息队列
		go queue.Run(cxt, wait)

	}

}

//所有子服务启动项函数
func boot(cxt context.Context, wait *sync.WaitGroup, httpOk chan bool, httpFail chan bool, sigs chan os.Signal) {

	defer func() {

		httpOk <- true

	}()

	//检查redis
	pingTimeoutCxt, c := context.WithTimeout(context.Background(), 1*time.Second)

	_, pingErr := redis.GetClient().Ping(pingTimeoutCxt).Result()

	c()

	if pingErr != nil {

		fmt.Println("redis连接失败，请检查")

		//发送信号让程序退出
		sigs <- syscall.SIGTERM

		return
	}

	client := http.Client().SetTimeout(1 * time.Second)

	for {

		select {

		//如http服务启动失败，其他子服务无需启动
		case <-httpFail:

			fmt.Println("http启动失败")

			//发送信号让程序退出
			sigs <- syscall.SIGTERM

			return

		default:

			time.Sleep(200 * time.Millisecond)

			//验证http服务已启动完成
			str, err := client.Request().GetToString("http://127.0.0.1:" + os.Getenv("PORT") + "/ping/" + kernel.Id)

			//http服务启动完成后再启动子服务
			if err == nil && str == "success" {

				//开启任务调度
				go crontab.Run(wait)

				//队列启动
				queueInit(cxt, wait)

				//记录pid和启动命令
				runInit()

				return

			}
		}

	}

}

//后台运行
func daemonize(args ...string) {

	//后台运行模式记录重定向输出

	sysType := runtime.GOOS

	if sysType == `windows` {

		cmd := gcmd2.NewCommand(tools.Join(" ", args)+" > logs/outLog.log", context.TODO())

		err := cmd.StartNoWait()

		if err != nil {

			log.Println(err)
		}

		return
	}

	if sysType == "linux" || sysType == "darwin" {

		runUser := os.Getenv("RUN_USER")

		if runUser == "" || runUser == "nobody" {

			cmd := gcmd2.NewCommand(tools.Join(" ", args)+" > logs/outLog.log 2>&1", context.TODO())

			err := cmd.StartNoWait()

			if err != nil {

				log.Println(err)
			}

			return

		}

		//以其他用户运行服务，源命令(sudo -u nginx ./main start)
		cmd := gcmd2.NewCommand("sudo -u "+runUser+" "+tools.Join(" ", args)+" > logs/outLog.log 2>&1", context.TODO())

		err := cmd.StartNoWait()

		if err != nil {

			log.Println(err)
		}

		return
	}

	fmt.Println("平台暂不支持")

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

		//以其他用户运行服务，源命令(sudo -u nginx ./main start)
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

//常规当前用户运行模式
func normal(args ...string) {

	cmd := gcmd2.NewCommand(tools.Join(" ", args), context.TODO())

	err := cmd.Start()

	if err != nil {

		log.Println(err)
	}

}

func stop() error {

	fmt.Println("stopping!!")

	b, err := common.PathExists("logs/run.pid")

	if err != nil {

		return err
	}

	if !b {

		return errors.New("run.pid文件不存在")
	}

	if b {

		pid, err := read.Open("logs/run.pid").Read()

		if err != nil {

			return err

		}

		sysType := runtime.GOOS

		var cmd *exec.Cmd

		if sysType == `windows` {

			//cmd = exec.Command("cmd", "/c", "taskkill /f /pid "+string(pid))
			cmd = exec.Command("cmd", "/c", ".\\lib\\windows-kill.exe -SIGINT "+string(pid))

		}

		if sysType == `linux` {

			cmd = exec.Command("bash", "-c", "kill "+string(pid))
		}

		err = cmd.Start()

		if err != nil {

			return err

		}

		err = cmd.Wait()

		if err != nil {

			return err

		}

		if sysType == `linux` {

			//等待进程退出
			for {

				time.Sleep(200 * time.Millisecond)

				wait := gcmd2.NewCommand("ps -p "+string(pid)+" | wc -l", context.TODO())

				num, waitErr := wait.CombinedOutput()

				str := strings.Replace(string(num), " ", "", -1)
				// 去除换行符
				str = strings.Replace(str, "\n", "", -1)

				if waitErr != nil {

					return waitErr

				}

				if str == "2" {

					continue

				}

				if str == "1" {

					fmt.Println("stopped!!")

					return nil
				}

			}

		}

		if sysType == `windows` {

			for {

				time.Sleep(200 * time.Millisecond)

				wait := gcmd2.NewCommand("tasklist|findstr   "+string(pid), context.TODO())

				_, waitErr := wait.CombinedOutput()

				if waitErr != nil {

					//signal.

					fmt.Println("stopped!!")

					return nil
				}

			}

		}

	}

	return nil
}

//记录pid
func runInit() {

	f, err := os.OpenFile("logs/run.pid", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)

	if err != nil {

		panic(err)
	}

	//记录pid
	_, err = f.Write([]byte(cast.ToString(os.Getpid())))

	if err == nil {

		isRun = true
	}

	_ = f.Close()

}
