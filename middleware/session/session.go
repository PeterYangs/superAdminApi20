package session

import (
	"context"
	"encoding/json"
	"gin-web/component/logs"
	"gin-web/conf"
	"gin-web/contextPlus"
	"gin-web/redis"
	"github.com/PeterYangs/tools"
	uuid "github.com/satori/go.uuid"
	"os"
	"strconv"
	"time"
)

func StartSession(c *contextPlus.Context) {

	//获取浏览器中的cookie
	cookie, err := c.Cookie(conf.Get("cookie_name").(string))

	//配置中读取session的存活时间
	sessionLifetime_, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	sessionLifetime := int64(sessionLifetime_)

	if err == nil {

		//查询redis中是否存在
		cookies, err := redis.GetClient().Get(context.TODO(), contextPlus.GetRedisSessionKey(cookie)).Result()

		if err == nil {

			//redis中存在的话，就将过期时间更新
			var session contextPlus.Session

			json.Unmarshal([]byte(cookies), &session)

			session.ExpireTime = time.Now().Unix() + sessionLifetime

			s, err := json.Marshal(session)

			if err != nil {

				logs.NewLogs().Error(err.Error())

				return
			}

			c.Set(conf.Get("cookie_key").(string), session)

			redis.GetClient().Set(context.TODO(), contextPlus.GetRedisSessionKey(session.Cookie), s, time.Second*time.Duration(sessionLifetime))

			return
		}

	}

	//不存在就设置session值
	u := uuid.NewV4().String() + "-" + tools.Md5(os.Getenv("APP_NAME"))

	life, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	c.SetCookie(conf.Get("cookie_name").(string), u, life, "/", c.Domain(), false, true)

	session := contextPlus.Session{
		Cookie:      u,
		ExpireTime:  time.Now().Unix() + sessionLifetime,
		SessionList: make(map[string]interface{}),
	}

	s, err := json.Marshal(session)

	if err != nil {

		logs.NewLogs().Error(err.Error())

		return
	}

	c.Set(conf.Get("cookie_key").(string), session)

	redis.GetClient().Set(context.TODO(), contextPlus.GetRedisSessionKey(u), s, time.Second*time.Duration(sessionLifetime))

}
