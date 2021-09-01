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
	"github.com/spf13/cast"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	kernel.IdInit()

	sigs := make(chan os.Signal, 1)

	//退出信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cxt, cancel := context.WithCancel(context.Background())

	wait := sync.WaitGroup{}

	httpOk := make(chan bool)

	httpFail := make(chan bool)

	go func() {

		sig := <-sigs

		fmt.Println()
		fmt.Println(sig)

		cancel()
	}()

	go boot(cxt, &wait, httpOk, httpFail)

	//启动http服务
	go httpStart(httpFail)

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
