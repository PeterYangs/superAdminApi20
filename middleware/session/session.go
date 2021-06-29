package session

import (
	"context"
	"encoding/json"
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

	cookie, err := c.Cookie(conf.Get("cookie_name").(string))

	if err == nil {

		cookies, err := redis.GetClient().Get(context.TODO(), contextPlus.GetRedisSessionKey(cookie)).Result()

		if err == nil {

			var session contextPlus.Session

			json.Unmarshal([]byte(cookies), &session)

			c.Set(conf.Get("cookie_key").(string), session)

			//c.Abort()

			return
		}

	}

	u := uuid.NewV4().String() + "-" + tools.Md5(os.Getenv("APP_NAME"))

	life, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	c.SetCookie(conf.Get("cookie_name").(string), u, life, "/", c.Domain(), false, true)

	sessionLifetime_, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	sessionLifetime := int64(sessionLifetime_)

	session := contextPlus.Session{
		Cookie:      u,
		ExpireTime:  time.Now().Unix() + sessionLifetime,
		SessionList: make(map[string]interface{}),
	}

	s, err := json.Marshal(session)

	c.Set(conf.Get("cookie_key").(string), session)

	if err != nil {

		return
	}

	redis.GetClient().Set(context.TODO(), contextPlus.GetRedisSessionKey(u), s, time.Second*time.Duration(sessionLifetime))

}
