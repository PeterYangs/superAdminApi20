package controller

import "github.com/gin-gonic/gin"

// Index 主页
func Index(c *gin.Context) interface{} {

	return gin.H{"code": 1, "msg": "hello world"}
}
