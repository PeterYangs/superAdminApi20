# superAdmin

开箱即用的后台框架

[中文文档](https://www.kancloud.cn/peter_yang/v001/2401726)

### 在线demo

https://www.peterdemo.net

账号：test

密码：Aa123456

### 前端仓库
https://github.com/PeterYangs/superAdminPage20

### 环境要求

redis

mysql

### Quick start

**controller**

```go
package controller

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
	"gin-web/response"
)

// Index 主页
func Index(c *contextPlus.Context) *response.Response {

	return response.Resp().Json(gin.H{"data": "hello world"})
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
func Session(c *contextPlus.Context) *response.Response {

    c.Session.Set("key", "value")

    c.Session.Get("key")

    return nil
}
```

**全局中间件**

```go
package kernel

import (
	"gin-web/contextPlus"
	"gin-web/middleware/accessLog"
	"gin-web/middleware/exception"
	"gin-web/middleware/session"
)

// Middleware 全局中间件
var Middleware []contextPlus.HandlerFunc

func Load() {

	Middleware = []contextPlus.HandlerFunc{
		exception.Exception,
		session.StartSession,
		accessLog.AccessLog,
	}

}
```



**验证码**

获取验证码

```go
func Captcha(c *contextPlus.Context) *response.Response {

    b := c.GetCaptcha()

    c.Header("content-type", "image/png")

    return response.Resp().Byte(b)
}
```

检查验证码

```go
func CheckCaptcha(c *contextPlus.Context) *response.Response {

    code := c.Query("code")

    bool:= c.CheckCaptcha(code)

    return response.Resp().Json(gin.H{"bool":code})

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
func Regex(c *contextPlus.Context) *response.Response {

	type regex struct {
		Test []string `form:"test[]" json:"test" regex:"[0-9a-z/]+"`
		Name string   `uri:"name" json:"name" regex:"[0-9a-z]+"`
	}

	var t regex

	err := c.ShouldBindPlus(&t)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	return response.Resp().Json(gin.H{"code": 1, "msg": "hello world"})
}

```

**数据库迁移**

```shell

[root@localhost ~]# go run  .\artisan\bin\artisan.go
Use the arrow keys to navigate: ↓ ↑ → ←
? 选择类型:
  > 数据库迁移
    数据填充
    生成key
    生成任务类

```

迁移文件

```go
package migrate_2019_08_12_055619_create_admin_table

import "gin-web/migrate"

func Up() {

	migrate.Create("admin", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2019_08_12_055619_create_admin_table"

		//主键
		createMigrate.BigIncrements("id")

		//int
		createMigrate.Integer("user_id").Unsigned().Nullable().Default(0).Unique().Comment("用户id")

		//varchar
		createMigrate.String("title", 255).Default("").Comment("标题")

		//text
		createMigrate.Text("content").Default(migrate.Null).Comment("内容")

		//索引
		createMigrate.Unique("user_id", "title")

	})

}

func Down() {

	migrate.DropIfExists("admin")

}

```

**限流器**

```go
package loginLimiter

import (
	"gin-web/component/limiter"
	"gin-web/contextPlus"
	"golang.org/x/time/rate"
	"time"
)

func LoginLimiter(c *contextPlus.Context) {

	//每秒生成一个令牌，桶的大小是10,第三个参数是自定义key，根据自定义的key寻找限流器（默认是每1分钟清理一次过期的限流器）
	if !limiter.NewLimiter(rate.Every(1*time.Second), 10, c.ClientIP()).Allow() {

		c.String(500, "访问频率过高")

		c.Abort()

	}

}

```

**分布式锁**

非阻塞

```go
func Index(c *contextPlus.Context) *response.Response {

    //申请一个锁，过期时间是10秒
    lock := redis.GetClient().Lock("lock", 10*time.Second)

    //释放锁
    defer lock.Release()

    //是否拿到锁
    if lock.Get() {

        return response.Resp().Json(gin.H{"res": true})
    
    }

    return response.Resp().Json(gin.H{"res": false})

}
```

阻塞

```go
func Index(c *contextPlus.Context) *response.Response {

    //申请一个锁，过期时间是10秒
    lock := redis.GetClient().Lock("lock", 10*time.Second)

    defer lock.Release()

    //是否拿到锁
    if lock.Block(time.Second * 3) {
    	
    	time.Sleep(4 * time.Second)

        return response.Resp().Json(gin.H{"res": true})
    }

    return response.Resp().Json(gin.H{"res": false})

}
```

**消息队列**

生成任务类
```shell
[root@localhost superAdminApi20]# go run artisan/bin/artisan.go
Use the arrow keys to navigate: ↓ ↑ → ←
? 选择类型:
    数据库迁移
    数据填充
    生成key
  > 生成任务类
```

任务类
```go
package access

import (
	"gin-web/database"
	"gin-web/model"
	"gin-web/task"
)

type TaskAccess struct {
	task.BaseTask
	Parameters *Parameter
}

type Parameter struct {
	task.Parameter
	Ip      string
	Url     string
	Params  string
	AdminId float64
}

func NewTask(ip string, url string, params string, adminId float64) *TaskAccess {

	return &TaskAccess{

		BaseTask: task.BaseTask{
			TaskName: "access",
		},
		Parameters: &Parameter{
			Ip:      ip,
			Url:     url,
			Params:  params,
			AdminId: adminId,
		},
	}
}

func (t *TaskAccess) Run() {

	database.GetDb().Create(&model.Access{
		Ip:      t.Parameters.Ip,
		Url:     t.Parameters.Url,
		Params:  t.Parameters.Params,
		AdminId: t.Parameters.AdminId,
	})

}

func (t *TaskAccess) BindParameters(p map[string]interface{}) {

	t.BaseTask.Bind(t.Parameters, p)

}

```


即时任务

```go
package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/response"
	"gin-web/task/email"
	"gin-web/task/sms"
)

func Task(c *contextPlus.Context) *response.Response {

	queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Run()

	return response.Resp().Api(1, "123", "")

}
```

延迟队列

```go
package controller

import (
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/response"
	"gin-web/task/email"
	"gin-web/task/sms"
	"time"
)

func Task(c *contextPlus.Context) *response.Response {

	queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).Queue("low").Delay(100 * time.Second).Run()

	return response.Resp().Api(1, "123", "")

}
```
**任务调度**

crontab/conf.go
```go
package crontab

import "fmt"

func Registered(c *crontab) {

	c.newSchedule().everyHour().function(func() {

		fmt.Println("每小时")

	})

	c.newSchedule().hourlyAt(16).everyMinute().function(func() {

		fmt.Println("每个16点的每分钟")

	})

	c.newSchedule().minuteAt(18).function(func() {

		fmt.Println("每小时的第18分钟")

	})

	c.newSchedule().everyMinute().function(func() {

		//panic("模拟报错")

		fmt.Println("每分钟")

	})

	c.newSchedule().everyMinuteAt(2).function(func() {

		fmt.Println("每2分钟")

	})

	c.newSchedule().everyDay().hourlyAt(16).minuteAt(36).function(func() {

		fmt.Println("每天16点36分")

	})

	c.newSchedule().dayAt(23).hourlyAt(16).minuteAt(50).function(func() {

		fmt.Println("23号16点50分")

	})

	c.newSchedule().dayAt(24).hourBetween(8, 10).function(func() {

		fmt.Println("24号8点-10点")

	})

	c.newSchedule().hourBetween(8, 9).everyMinute().function(func() {

		fmt.Println("24号8点-9点每分钟")

	})

	c.newSchedule().dayBetween(22, 24).everyHour().everyMinute().function(func() {

		fmt.Println("22号-24号每分钟")

	})

}

```






