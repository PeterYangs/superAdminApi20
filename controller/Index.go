package controller

import "github.com/gin-gonic/gin"

//主页
func Index(c *gin.Context) gin.H {

	return gin.H{"code": 1, "msg": "hello world"}
}
