package routes

import (
	"gin-web/controller"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

const (
	GET    int = 0x000000
	POST   int = 0x000001
	PUT    int = 0x000002
	DELETE int = 0x000003
)

func Load(rr *gin.Engine) {

	//绑定到全局变量
	r = rr

	actionRegistered(GET, "/", controller.Index, middleware.GoOn)
	//actionRegistered(GET, "/query", controller.Query)

}

//路由注册
func actionRegistered(method int, url string, f func(c *controller.Contexts) interface{}, middlewares ...gin.HandlerFunc) {

	ff := func(c *gin.Context) {

		data := f(&controller.Contexts{c})

		getDataType(data, &controller.Contexts{c})

	}

	middlewares = append(middlewares, ff)

	switch method {

	case GET:

		r.GET(url, middlewares...)

	case POST:

		r.POST(url, middlewares...)

	case PUT:

		r.PUT(url, middlewares...)

	case DELETE:

		r.DELETE(url, middlewares...)

	}

}

func getDataType(data interface{}, c *controller.Contexts) {

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
