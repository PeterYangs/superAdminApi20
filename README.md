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
	"gin-web/controller/file"
	"gin-web/controller/regex"
	"gin-web/kernel"
	"gin-web/middleware"
	"github.com/gin-gonic/gin"
)

func _init(_r group) {

	//路由组，支持嵌套
	_r.Group("/api", func(g group) {

		g.Registered(GET, "/", controller.Index).Bind()
		g.Registered(GET, "/gg", controller.Index).Bind()

		g.Group("/login", func(g2 group) {

			g2.Registered(GET, "/", controller.Index).Bind()
		})

	}, middleware.GoOn)

	//单路由
	_r.Registered(GET, "/", controller.Index).Bind()
	
	_r.Registered(GET, "/check", controller.CheckCaptcha).Bind()



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

**参数验证**
```go
package regex

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
)

// Regex 参数规则验证示例，路由为 /regex/:name ,请求为 /regex/1sds?test[]=1&test[]=2,regex标记只支持string和[]string两个类型
func Regex(c *contextPlus.Context) interface{} {

	type regex struct {
		Test []string `form:"test[]" json:"test" regex:"[0-9a-z/]+"`
		Name string   `uri:"name" json:"name" regex:"[0-9a-z]+"`
	}

	var t regex

	err := c.ShouldBindPlus(&t)

	if err != nil {

		return gin.H{"code": 2, "mgs": err.Error()}

	}

	return gin.H{"code": 1, "msg": "hello world"}
}

```

**数据库迁移**

```shell
[root@localhost ~]# go run migrate/bin/run.go
Use the arrow keys to navigate: ↓ ↑ → ←
? 选择类型:
  > 创建数据库迁移
    执行迁移
    回滚迁移



```



迁移文件
```go
package migrate_2019_08_12_055619_create_admin_table

import "gin-web/migrate"

func Up() {

	migrate.Create("admin", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2019_08_12_055619_create_admin_table"

		createMigrate.BigIncrements("id")

		createMigrate.Integer("user_id").Unsigned().Nullable().Default(0).Comment("用户id")

		createMigrate.String("title", 255).Default("").Comment("标题")

	})

}

func Down() {

	migrate.DropIfExists("admin")

}


```


