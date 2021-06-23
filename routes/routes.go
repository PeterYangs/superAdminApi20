package routes

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
	"sync"
)

//var r *gin.Engine

const (
	GET    int = 0x000000
	POST   int = 0x000001
	PUT    int = 0x000002
	DELETE int = 0x000003
)

type router struct {
	engine *gin.Engine
	//regex  map[string]string //路由正则表达式
}

type group struct {
	group *gin.RouterGroup
	//regex map[string]string //路由正则表达式
}

type handler struct {
	handlerFunc []gin.HandlerFunc
	group       *gin.RouterGroup
	url         string
	method      int
}

func newRouter(engine *gin.Engine) *router {

	return &router{
		engine: engine,
	}
}

func (rr *router) Group(path string, callback func(group2 *group), middlewares ...contextPlus.HandlerFunc) {

	var temp = make([]gin.HandlerFunc, len(middlewares))

	//fmt.Println(middlewares)

	for i, funcs := range middlewares {

		//fmt.Println(funcs)

		tempFuncs := funcs

		temp[i] = func(context *gin.Context) {

			tempFuncs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

		//f := func(context *gin.Context) {
		//
		//	tempFuncs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})
		//
		//}

		//temp = append(temp, f)

	}

	//fmt.Println(temp)

	g := group{
		group: rr.engine.Group(path, temp...),
	}

	callback(&g)

}

func (gg *group) Group(path string, callback func(group2 *group), middlewares ...contextPlus.HandlerFunc) {

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		tempFuncs := funcs

		temp[i] = func(context *gin.Context) {

			tempFuncs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

	}

	g := group{
		group: gg.group.Group(path, temp...),
	}

	callback(&g)

}

//func (rr *router) Regex(r map[string]string) *router {
//
//	rr.regex = r
//
//	return rr
//}
//
//func (gg *group) Regex(r map[string]string) *group {
//
//	gg.regex = r
//
//	return gg
//}

func (rr *router) Registered(method int, url string, f func(c *contextPlus.Context) interface{}, middlewares ...contextPlus.HandlerFunc) {

	ff := func(c *contextPlus.Context) {

		data := f(c)

		getDataType(data, c)

	}

	middlewares = append(middlewares, ff)

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		tempFuncs := funcs

		temp[i] = func(context *gin.Context) {

			tempFuncs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

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

func (gg *group) Registered(method int, url string, f func(c *contextPlus.Context) interface{}, middlewares ...contextPlus.HandlerFunc) *handler {

	ff := func(c *contextPlus.Context) {

		data := f(c)

		getDataType(data, c)

	}

	middlewares = append(middlewares, ff)

	var temp = make([]gin.HandlerFunc, len(middlewares))

	for i, funcs := range middlewares {

		tempFuncs := funcs

		temp[i] = func(context *gin.Context) {

			tempFuncs(&contextPlus.Context{Context: context, Lock: &sync.Mutex{}})

		}

	}

	return &handler{
		handlerFunc: temp,
		group:       gg.group,
		url:         url,
		method:      method,
	}

}

func (h *handler) Bind() {

	switch h.method {

	case GET:

		h.group.GET(h.url, h.handlerFunc...)

	case POST:

		h.group.POST(h.url, h.handlerFunc...)

	case PUT:

		h.group.PUT(h.url, h.handlerFunc...)

	case DELETE:

		h.group.DELETE(h.url, h.handlerFunc...)
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
