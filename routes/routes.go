package routes

import (
	"gin-web/contextPlus"
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

type router struct {
	engine *gin.Engine
}

type group struct {
	group *gin.RouterGroup
}

func newRouter(engine *gin.Engine) *router {

	return &router{
		engine: engine,
	}
}

func (rr *router) Group(path string, callback func(group2 *group), middlewares ...contextPlus.HandlerFunc) {

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		temp[i] = func(context *gin.Context) {

			funcs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

	}

	g := group{
		group: rr.engine.Group(path, temp...),
	}

	callback(&g)

}

func (gg *group) Group(path string, callback func(group2 *group), middlewares ...contextPlus.HandlerFunc) {

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		temp[i] = func(context *gin.Context) {

			funcs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

	}

	g := group{
		group: gg.group.Group(path, temp...),
	}

	callback(&g)

}

func (rr *router) Registered(method int, url string, f func(c *contextPlus.Context) interface{}, middlewares ...contextPlus.HandlerFunc) {

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

		rr.engine.GET(url, temp...)

	case POST:

		rr.engine.POST(url, temp...)

	case PUT:

		rr.engine.PUT(url, temp...)

	case DELETE:

		rr.engine.DELETE(url, temp...)

	}

}

func (gg *group) Registered(method int, url string, f func(c *contextPlus.Context) interface{}, middlewares ...contextPlus.HandlerFunc) {

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

		gg.group.GET(url, temp...)

	case POST:

		gg.group.POST(url, temp...)

	case PUT:

		gg.group.PUT(url, temp...)

	case DELETE:

		gg.group.DELETE(url, temp...)

	}

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

	case []byte:

		c.String(200, string(item))
	case gin.H:

		//fmt.Println(1111)
		c.JSON(200, item)

	}
}
