package middleware

import (
	"fmt"
	"gin-web/contextPlus"
)

func GoOn(c *contextPlus.Context) {

	fmt.Println("after")

	c.Next()

	//这里为请求后处理
	fmt.Println("last")

}

func Stop(c *contextPlus.Context) {

	c.Abort()
}
