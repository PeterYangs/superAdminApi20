package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/response"
	"gin-web/task/email"
	"gin-web/task/sms"
	"time"
)

func Task(c *contextPlus.Context) *response.Response {

	queue.Dispatch(email.NewTask("1111111111111111111", "title", "content")).Queue("low").Delay(100 * time.Second).Run()
	//queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()

	return response.Resp().Api(1, "123", "")

}

func Task2(c *contextPlus.Context) *response.Response {

	queue.Dispatch(sms.NewTask("110", "content")).Run()

	return response.Resp().Api(1, "123", "13")

}
