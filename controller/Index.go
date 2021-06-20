package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Index 主页
func Index(c *gin.Context) interface{} {

	//fmt.Println("index")

	s, err := c.Cookie("gin_web")

	fmt.Println(s)
	fmt.Println(err)

	return gin.H{"code": 1, "msg": "hello world"}
}
