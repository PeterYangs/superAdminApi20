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


### 开发模式
```shell
go run main.go
```

### 编译&部署
编译二进制文件
```
go build main.go
```

直接运行
```
./main

或者

./main start
```

后台运行
```
./main start -d
```

安全停止
```
./main stop
```

平滑重启
```
./main restart
```

运行内置命令
```shell
./main artisan
```




### Quick start

**controller**

```go
package controller

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
)

func Index(c *contextPlus.Context) *response.Response {

	return response.Resp().Api(1, "success", "index")
}
```

**route**

routes/routes.go

```go
func Routes(_r route.Group) {

    _r.Registered(route.GET, "/index", controller.Index).Bind()

    _r.Group("/login", func (_login route.Group) {

        _login.Registered(route.POST, "/login", login.Login, loginLimiter.LoginLimiter).Bind()

        _login.Registered(route.POST, "/logout", login.Logout).Bind()

    })
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

middleware/global.go
```go
package middleware

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/middleware/session"
	"superadmin/middleware/accessLog"
)

func Load() []contextPlus.HandlerFunc {

	return []contextPlus.HandlerFunc{

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

[root@localhost ~]# ./main artisan
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
[root@localhost superAdminApi20]# ./main artisan
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

func (t *TaskAccess) Run() error {

	database.GetDb().Create(&model.Access{
		Ip:      t.Parameters.Ip,
		Url:     t.Parameters.Url,
		Params:  t.Parameters.Params,
		AdminId: t.Parameters.AdminId,
	})
	
	return nil

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


重试次数
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

	queue.Dispatch(email.NewTask("904801074@qq.com", "title", "content")).SetTryTime(3).Run()

	return response.Resp().Api(1, "123", "")

}
```


**任务调度**

crontab/crontabs.go
```go
package crontab

import "fmt"

func Crontab(c *crontab) {

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

		fmt.Println("每天8点-9点每分钟")

	})

	c.newSchedule().dayBetween(22, 24).everyHour().everyMinute().function(func() {

		fmt.Println("22号-24号每分钟")

	})

}

```


**命令行**

需要实现
```go
type Artisan interface {
	ArtisanRun()
	GetName() string
}
```
一个例子
```go
package demo

import "github.com/PeterYangs/superAdminCore/component/logs"

type Demo struct {
}

func (d Demo) GetName() string {

	return "demo"
}

func (d Demo) ArtisanRun() {

	logs.NewLogs().Debug("demo")
}

```
运行
```shell
./main artisan
```
**缓存**

缓存支持两种驱动`file`和`redis`，在`.env`文件中设置`CACHE_DRIVER=redis`
```go
package test

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/cache"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
)

func Cache(c *contextPlus.Context) *response.Response {

	cache.Cache().Put("nice", "0", 0)

	fmt.Println(cache.Cache().Get("nice"))

	return response.Resp().Api(1, "success", "")
}

```




