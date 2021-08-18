package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/response"
	"gin-web/task/email"
)

func Task(c *contextPlus.Context) *response.Response {

	//queue.Dispatch(email.NewTask(cast.ToString(tools.MtRand(10, 1000)), "title", "content")).Delay(100 * time.Hour).Queue("low").Run()
	queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()

	return response.Resp().Api(1, "123", "")

}

func Task2(c *contextPlus.Context) *response.Response {

	//queue.Dispatch(sms.NewTask("110", "content")).Run()

	return response.Resp().Api(1, "123", "13")

}
