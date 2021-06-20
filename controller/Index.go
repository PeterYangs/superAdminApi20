package controller

import (
	"context"
	"fmt"
	"gin-web/redis"
	"github.com/gin-gonic/gin"
	"time"
)

type Contexts struct {
	*gin.Context
}

type HandlerFuncs struct {
	gin.HandlerFunc
}

func (c *Contexts) Domain() string {
	//return nameOfFunction(c.handlers.Last())

	return "11111111111111111111"
}

// Index 主页
func Index(c *Contexts) interface{} {

	//fmt.Println("index")

	//s, err := c.Cookie("gin_web")

	//fmt.Println(s)
	//fmt.Println(err)

	//c.Domain()
	//c.Request

	fmt.Println(c.Domain())

	//fmt.Println(c.Request.)

	redis.GetClient().Set(context.TODO(), "gg", "yy", time.Second*1000)

	return gin.H{"code": 1, "msg": "hello world"}
}
