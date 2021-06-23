package kernel

import (
	"gin-web/contextPlus"
	"gin-web/middleware/exception"
	"gin-web/middleware/session"
	"github.com/gin-gonic/gin"
)

// Middleware 全局中间件
var Middleware []func(*contextPlus.Context)

func Load(r *gin.Engine) {

	Middleware = []func(*contextPlus.Context){
		exception.Exception,
		session.StartSession,
	}

	put(r)

}

func put(r *gin.Engine) {

	for _, f := range Middleware {

		globalMiddleware(r, f)
	}

}

// GlobalMiddleware 设置全局中间件
func globalMiddleware(r *gin.Engine, middleware func(*contextPlus.Context)) {

	//函数转换
	temp := func(context *gin.Context) {

		middleware(&contextPlus.Context{Context: context})

	}

	r.Use(temp)

}
