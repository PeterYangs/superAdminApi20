package main

import (
	"context"
	"fmt"
	"gin-web/component/logs"
	"gin-web/conf"
	"gin-web/crontab"
	"gin-web/kernel"
	"gin-web/queue"
	"gin-web/routes"
	"github.com/PeterYangs/tools/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/spf13/cast"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type program struct {
	sigs chan os.Signal

	cxt context.Context

	cancel context.CancelFunc

	wait *sync.WaitGroup

	httpOk chan bool

	httpFail chan bool
}

func (p *program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {
	//log.Println("停止服务")

	//time.Sleep(1*time.Second)

	//if runtime.GOOS == "linux" {
	//
	//	pid, err := read.Open("logs/run.pid").Read()
	//
	//	if err != nil {
	//
	//		return err
	//	}
	//
	//	_, err = gcmd.Command("kill " + string(pid)).Start()
	//
	//	if err != nil {
	//
	//		return err
	//	}
	//
	//}

	<-p.httpOk

	//等待队列退出
	p.wait.Wait()

	return nil
}
func (p *program) run() {
	// 这里放置程序要执行的代码……

	//fmt.Println("gg")

	p.serviceStart()
}

func main() {

	kernel.IdInit()

	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	httpOk := make(chan bool)

	httpFail := make(chan bool)

	//go boot(cxt, &wait, httpOk, httpFail)

	//启动http服务
	//go httpStart(httpFail)

	cfg := &service.Config{
		Name:        "simple",
		DisplayName: "a simple service",
		Description: "This is an example Go service.",
	}
	// Interface 接口
	prg := &program{

		sigs:     sigs,
		cxt:      cxt,
		cancel:   cancel,
		httpOk:   httpOk,
		httpFail: httpFail,
		wait:     &wait,
	}
	// 构建服务对象
	s, err := service.New(prg, cfg)
	if err != nil {
		log.Fatal(err)
	}
	// logger 用于记录系统日志
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 2 { //如果有命令则执行
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else { //否则说明是方法启动了
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
	}
	if err != nil {
		logger.Error(err)
	}

}

func (p *program) serviceStart() {

	//kernel.IdInit()
	//
	//sigs := make(chan os.Signal, 1)
	//
	////退出信号
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	//cxt, cancel := context.WithCancel(context.Background())
	//
	//wait := sync.WaitGroup{}
	//
	//httpOk := make(chan bool)
	//
	//httpFail := make(chan bool)
	//
	//go boot(cxt, &wait, httpOk, httpFail)
	//
	////启动http服务
	//go httpStart(httpFail)

	go func() {

		sig := <-p.sigs

		fmt.Println()
		fmt.Println(sig)

		p.cancel()
	}()

	go boot(p.cxt, p.wait, p.httpOk, p.httpFail)

	//启动http服务
	go httpStart(p.httpFail)

	//<-p.httpOk

	//等待队列退出
	//p.wait.Wait()

	//fmt.Println("finish")

	p.wait.Wait()

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

func httpStart(httpFail chan bool) {

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

	err := r.Run(":" + port)

	if err != nil {

		log.Println(err)

		httpFail <- true
	}

}

func queueInit(cxt context.Context, wait *sync.WaitGroup) {

	//fmt.Println("哈哈哈")

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

			//fmt.Println(str, err)

			if err == nil && str == "success" {

				//开启任务调度
				go crontab.Run(wait)

				//队列启动
				queueInit(cxt, wait)

				//日志模块初始化
				logInit(cxt, wait)

				return

			}
		}

	}

}
