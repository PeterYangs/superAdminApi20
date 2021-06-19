package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GoOn(c *gin.Context) {

	fmt.Println("after")

	c.Next()

	//这里为请求后处理
	fmt.Println("last")

}

func Stop(c *gin.Context) {

	c.Abort()
}
