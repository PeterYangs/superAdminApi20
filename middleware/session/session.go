package session

import (
	"context"
	"encoding/json"
	"gin-web/contextPlus"
	"gin-web/redis"
	"github.com/PeterYangs/tools"
	uuid "github.com/satori/go.uuid"
	"os"
	"strconv"
	"time"
)

//type Session struct {
//	Cookie      string                 `json:"cookie"`
//	ExpireTime  int64                  `json:"expire_time"`
//	SessionList map[string]interface{} `json:"session_list"`
//}

func StartSession(c *contextPlus.Context) {

	defer c.Next()

	cookie, err := c.Cookie(os.Getenv("APP_NAME") + "_session")

	if err == nil {

		cookies, err := redis.GetClient().Get(context.TODO(), "session:"+cookie).Result()

		if err == nil {

			//fmt.Println("find")

			//c.CookieKey=cookie

			var session contextPlus.Session

			json.Unmarshal([]byte(cookies), &session)

			c.Set("cookie_key", session)

			return
		}

	}

	u := uuid.NewV4().String() + "-" + tools.Md5(os.Getenv("APP_NAME"))

	life, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	c.SetCookie(os.Getenv("APP_NAME")+"_session", u, life, "/", c.Domain(), false, true)

	sessionLifetime_, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	sessionLifetime := int64(sessionLifetime_)

	session := contextPlus.Session{
		Cookie:      u,
		ExpireTime:  time.Now().Unix() + sessionLifetime,
		SessionList: make(map[string]interface{}),
	}

	s, err := json.Marshal(session)

	c.Set("cookie_key", session)

	if err != nil {

		return
	}

	redis.GetClient().Set(context.TODO(), "session:"+u, s, time.Second*time.Duration(sessionLifetime))

}
