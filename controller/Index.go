package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Index 主页
func Index(c *gin.Context) interface{} {

	fmt.Println("index")

	return gin.H{"code": 1, "msg": "hello world"}
}
