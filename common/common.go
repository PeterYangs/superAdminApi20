package common

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
)

// GlobalMiddleware 设置全局中间件
func GlobalMiddleware(r *gin.Engine, middleware func(*contextPlus.Context)) {

	//函数转换
	temp := func(context *gin.Context) {

		middleware(&contextPlus.Context{Context: context})

	}

	r.Use(temp)

}
