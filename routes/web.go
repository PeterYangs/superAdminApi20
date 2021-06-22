package routes

import (
	"gin-web/controller"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func Load(rr *gin.Engine) {

	_r := newRouter(rr)

	//路由组，支持嵌套
	_r.Group("/api", func(g *group) {

		g.Registered(GET, "/", controller.Index)
		g.Registered(GET, "/gg", controller.Index)

		g.Group("login", func(g2 *group) {

			g2.Registered(GET, "/", controller.Index)
		})

	}, middleware.GoOn)

	//单路由
	_r.Registered(GET, "/", controller.Index)

	_r.Registered(GET, "/c", controller.Captcha)

}
