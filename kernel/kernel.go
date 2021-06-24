package kernel

import (
	"gin-web/contextPlus"
	"gin-web/middleware/exception"
	"gin-web/middleware/session"
)

// Middleware 全局中间件
var Middleware []contextPlus.HandlerFunc

func Load() {

	Middleware = []contextPlus.HandlerFunc{
		exception.Exception,
		session.StartSession,
		//middleware.GoOn,
		//routeRegex.RouterRegex,
	}

}
