package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/redis"
	"gin-web/response"
	"gin-web/task/sms"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"golang.org/x/net/context"
	"time"
)

func Task(c *contextPlus.Context) *response.Response {

	//for i := 0; i < 100; i++ {
	//
	//	queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()
	//
	//	queue.Dispatch(sms.NewTask("110", "content")).Queue("high").Run()
	//}

	//queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Delay(100*time.Second).Run()
	//queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()

	re, _ := redis.GetClient().ZRangeByScore(context.TODO(), "queue:delay", &redis2.ZRangeBy{Min: "", Max: cast.ToString(time.Now().Unix())}).Result()

	return response.Resp().Api(1, "123", re)

}

func Task2(c *contextPlus.Context) *response.Response {

	queue.Dispatch(sms.NewTask("110", "content")).Run()

	return response.Resp().Api(1, "123", "13")

}
