package controller

import (
	"gin-web/contextPlus"
	"gin-web/redis"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

//var m =make(map[string]string)

// Index 主页
func Index(c *contextPlus.Context) *response.Response {

	//申请一个锁，过期时间是10秒
	lock := redis.GetClient().Lock("lock", 1*time.Second)

	defer lock.Release()

	//是否拿到锁
	if lock.Get() {

		time.Sleep(4 * time.Second)

		return response.Resp().Json(gin.H{"res": true})
	}

	return response.Resp().Json(gin.H{"res": false})

}

func Block(c *contextPlus.Context) *response.Response {

	//申请一个锁，过期时间是10秒
	lock := redis.GetClient().Lock("lock", 10*time.Second)

	defer lock.Release()

	//最多等待5秒
	if lock.Block(5 * time.Second) {

		time.Sleep(100 * time.Millisecond)

		return response.Resp().Json(gin.H{"res": true})

	}

	return response.Resp().Json(gin.H{"res": false})
}

//func Block2(c *contextPlus.Context) *response.Response {
//
//	//申请一个锁，过期时间是10秒
//	lock := redis.GetClient().Lock("lock", 10*time.Second)
//
//	defer lock.Release()
//
//	//最多等待5秒
//	if lock.Block(5 * time.Second) {
//
//		//time.Sleep(1 * time.Minute)
//
//		return response.Resp().Json(gin.H{"res": true})
//
//	}
//
//	return response.Resp().Json(gin.H{"res": false})
//}

func Index2(c *contextPlus.Context) *response.Response {

	lock := redis.GetClient().Lock("lock", 1000*time.Second)

	defer lock.Release()

	if lock.Get() {

		return response.Resp().Json(gin.H{"res": true})
	}

	//m["1"]="2"

	return response.Resp().Json(gin.H{"res": false})

	//return response.Resp().Json(gin.H{"code": 1, "msg": "hello world"})
}

// SessionSet 并发写入demo
func SessionSet(c *contextPlus.Context) *response.Response {

	s := c.Session()

	for i := 0; i < 100; i++ {

		go func(ii int, ss *contextPlus.Session) {

			ss.Set("test"+strconv.Itoa(ii), ii)

		}(i, s)
	}

	return response.Resp().Json(gin.H{"code": 1, "msg": "hello world"})
}

func Captcha(c *contextPlus.Context) *response.Response {

	b := c.GetCaptcha()

	c.Header("content-type", "image/png")

	return response.Resp().Byte(b)
}

func CheckCaptcha(c *contextPlus.Context) *response.Response {

	code := c.Query("code")

	return response.Resp().Json(gin.H{"bool": code})

}
