package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/response"
	"gin-web/task/email"
	"github.com/PeterYangs/tools"
	"github.com/spf13/cast"
	"time"
)

func Task(c *contextPlus.Context) *response.Response {

	queue.Dispatch(email.NewTask(cast.ToString(tools.MtRand(10, 1000)), "title", "content")).Delay(30 * time.Second).Queue("low").Run()
	//queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()

	//fmt.Println(123)

	//lock := redis.GetClient().Lock("test", 3*time.Second)
	//
	//defer lock.Release()
	//
	//if lock.Block(3 * time.Second) {
	//
	//	time.Sleep(5 * time.Second)
	//
	//	fmt.Println("拿锁成功")
	//
	//	return response.Resp().Api(1, "123", "拿锁成功")
	//}

	return response.Resp().Api(1, "123", "fail")

}

func Task2(c *contextPlus.Context) *response.Response {

	//queue.Dispatch(sms.NewTask("110", "content")).Run()

	return response.Resp().Api(1, "123", "13")

}
