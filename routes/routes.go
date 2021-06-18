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
func actionRegistered(method string, url string, f func(c *gin.Context) interface{}) {

	switch method {

	case "get":

		r.GET(url, func(c *gin.Context) {

			data := f(c)

			getDataType(data, c)

		})

	}

}

func getDataType(data interface{}, c *gin.Context) {

	switch item := data.(type) {

	case map[string]interface{}:

		c.JSON(200, item)

	case string:

		c.String(200, item)
	case gin.H:

		//fmt.Println(1111)
		c.JSON(200, item)

	}
}
