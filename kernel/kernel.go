package kernel

import (
	"gin-web/contextPlus"
	"gin-web/middleware/accessLog"
	"gin-web/middleware/exception"
	"gin-web/middleware/session"
	uuid "github.com/satori/go.uuid"
)

// Middleware 全局中间件
var Middleware []contextPlus.HandlerFunc

// Id 服务id
var Id string

func Load() {

	Middleware = []contextPlus.HandlerFunc{
		exception.Exception,
		session.StartSession,
		accessLog.AccessLog,
	}

}

func IdInit() {

	Id = uuid.NewV4().String()
}
