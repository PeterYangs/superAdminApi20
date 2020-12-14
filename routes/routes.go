package routes

import (
	"gin-web/controller"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Load(rr *gin.Engine) {

	//绑定到全局变量
	r = rr

	actionRegistered("get", "/", controller.Index)
	actionRegistered("get", "/query", controller.Query)

}

//路由注册
func actionRegistered(method string, url string, f func(c *gin.Context) gin.H) {

	switch method {

	case "get":

		r.GET(url, func(c *gin.Context) {

			data := f(c)

			c.JSON(200, data)

		})

	}

}
