# gin-web

基于gin封装的mc框架

### 环境要求

redis

### Quick start

**controller**

```go
package controller

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
)

// Index 主页
func Index(c *contextPlus.Context) interface{} {

	return gin.H{"code": 1, "msg": "hello world"}
}
```

**route**

route/web.go
```go
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

}
```

**session**
```go
func Session(c *contextPlus.Context) interface{} {
	
    c.Session.Set("key","value")
	
    c.Session.Get("key")
	
    return nil
}
```

**验证码**

获取验证码

```go
func Captcha(c *contextPlus.Context) interface{} {
	
    b:=c.GetCaptcha()

    c.Header("content-type", "image/png")

    return b
}
```

检查验证码
```go
func CheckCaptcha(c *contextPlus.Context) interface{} {

    code := c.Query("code")

	bool:=c.CheckCaptcha(code)
    
    return gin.H{"data":bool }

}
```

