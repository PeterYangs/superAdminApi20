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
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	//启动http服务
	go start()

	//延迟队列的标记
	wait.Add(1)

	for i := 0; i < cast.ToInt(os.Getenv("QUEUE_NUM")); i++ {

		wait.Add(1)

		//启动消息队列
		go queue.Run(cxt, &wait)

	}

	//日志模块初始化
	logInit(cxt, &wait)

	go func() {

		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)

		cancel()
	}()

	//等待队列退出
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

func start() {

	r := gin.Default()

	//加载配置
	conf.Load()

	//加载全局中间件
	kernel.Load()

	//加载路由
	routes.Load(r)

	//开启任务调度
	go crontab.Run()

	//设置端口
	port := os.Getenv("PORT")

	if port == "" {

		port = "8887"
	}

	r.Run(":" + port)

}
