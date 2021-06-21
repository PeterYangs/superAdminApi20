package contextPlus

import (
	"context"
	"encoding/json"
	"errors"
	"gin-web/conf"
	"gin-web/redis"
	"sync"
	"time"

	//"gin-web/structure"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
	//CookieKey string
	Lock *sync.Mutex
}

type Session struct {
	Cookie      string                 `json:"cookie"`
	ExpireTime  int64                  `json:"expire_time"`
	SessionList map[string]interface{} `json:"session_list"`
	Lock        *sync.Mutex
}

type HandlerFunc func(*Context)

func (c *Context) Domain() string {

	//c.Mu

	return tools.Explode(":", c.Request.Host)[0]
}

func (c *Context) Session() *Session {

	var session Session

	cookie_, ex := c.Get(conf.Get("cookie_key").(string))

	if ex {

		session = cookie_.(Session)
	}

	session.Lock = &sync.Mutex{}

	return &session
}

func (s *Session) Set(key string, value interface{}) error {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.GetClient().Get(context.TODO(), "session:"+s.Cookie).Result()

	if err != nil {

		return err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	//fmt.Println(session)

	if err != nil {

		return err
	}

	session.SessionList[key] = value

	sessionStringNew, err := json.Marshal(session)

	redis.GetClient().Set(context.TODO(), "session:"+s.Cookie, sessionStringNew, time.Duration(s.ExpireTime-time.Now().Unix())*time.Second)

	return nil
}

func (s *Session) Get(key string) (interface{}, error) {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.GetClient().Get(context.TODO(), "session:"+s.Cookie).Result()

	if err != nil {

		return nil, nil
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	value, ok := session.SessionList[key]

	if ok {

		return value, nil
	}

	return nil, errors.New("not found key is " + key)

}
