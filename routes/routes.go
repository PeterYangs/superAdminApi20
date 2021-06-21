package routes

import (
	"gin-web/contextPlus"
	"gin-web/controller"
	"github.com/gin-gonic/gin"
	"sync"
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

	actionRegistered(GET, "/", controller.Index)
	//actionRegistered(GET, "/query", controller.Query)

}

//路由注册
func actionRegistered(method int, url string, f func(c *contextPlus.Context) interface{}, middlewares ...contextPlus.HandlerFunc) {

	ff := func(c *contextPlus.Context) {

		data := f(c)

		getDataType(data, c)

	}

	middlewares = append(middlewares, ff)

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		temp[i] = func(context *gin.Context) {

			funcs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

	}

	switch method {

	case GET:

		r.GET(url, temp...)

	case POST:

		r.POST(url, temp...)

	case PUT:

		r.PUT(url, temp...)

	case DELETE:

		r.DELETE(url, temp...)

	}

}

func getDataType(data interface{}, c *contextPlus.Context) {

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
