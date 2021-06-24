package middleware

import (
	"gin-web/contextPlus"
)

func GoOn(c *contextPlus.Context) {

	//fmt.Println("before")
	//
	//c.Next()
	//
	////这里为请求后处理
	//fmt.Println("after")

	//fmt.Println(c.Test)

	//fmt.Println(c.Regex)

}

func Stop(c *contextPlus.Context) {

	c.Abort()
}
