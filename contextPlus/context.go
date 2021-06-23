package contextPlus

import (
	"context"
	"encoding/json"
	"errors"
	"gin-web/component/captcha"
	"gin-web/conf"
	"gin-web/redis"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
	"sync"
	"time"
)

type Context struct {
	*gin.Context
	Lock *sync.Mutex
	//Regex string
	Regex map[string]string
	//Jj *map[string]string

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

// GetCaptcha 获取验证码
func (c *Context) GetCaptcha() []byte {

	text := captcha.GetRandStr(4)

	expiredTime := int(time.Now().Unix()) + cast.ToInt(conf.Get("captcha_lifetime"))

	value := text + "," + cast.ToString(expiredTime)

	c.Session().Set(conf.Get("captcha_key").(string), value)

	return captcha.ImgText(150, 60, text)

}

// CheckCaptcha 验证验证码
func (c *Context) CheckCaptcha(code string) bool {

	realCaptcha, err := c.Session().Get(conf.Get("captcha_key").(string))

	if err != nil {

		return false
	}
	//无论正确或错误，检查完成后移除这个验证码
	defer c.Session().Remove(conf.Get("captcha_key").(string))

	temp := tools.Explode(",", realCaptcha.(string))

	realCode := temp[0]

	expiredTime := temp[1]

	if strings.ToLower(realCode) == strings.ToLower(code) && cast.ToInt64(expiredTime) > time.Now().Unix() {

		return true
	}

	return false
}

type Session struct {
	Cookie      string                 `json:"cookie"`
	ExpireTime  int64                  `json:"expire_time"`
	SessionList map[string]interface{} `json:"session_list"`
	Lock        *sync.Mutex
}

func (s *Session) Set(key string, value interface{}) error {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.GetClient().Get(context.TODO(), GetRedisSessionKey(s.Cookie)).Result()

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

	redis.GetClient().Set(context.TODO(), GetRedisSessionKey(s.Cookie), sessionStringNew, time.Duration(s.ExpireTime-time.Now().Unix())*time.Second)

	return nil
}

func (s *Session) Get(key string) (interface{}, error) {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.GetClient().Get(context.TODO(), GetRedisSessionKey(s.Cookie)).Result()

	if err != nil {

		return nil, err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	if err != nil {

		return nil, err
	}

	value, ok := session.SessionList[key]

	if ok {

		return value, nil
	}

	return nil, errors.New("not found key is " + key)

}

func (s *Session) Remove(key string) error {

	s.Lock.Lock()

	defer s.Lock.Unlock()

	sessionString, err := redis.GetClient().Get(context.TODO(), GetRedisSessionKey(s.Cookie)).Result()

	if err != nil {

		return err
	}

	var session Session

	err = json.Unmarshal([]byte(sessionString), &session)

	if err != nil {

		return err
	}

	delete(session.SessionList, key)

	sessionStringNew, err := json.Marshal(session)

	if err != nil {

		return err
	}

	redis.GetClient().Set(context.TODO(), GetRedisSessionKey(s.Cookie), sessionStringNew, time.Duration(s.ExpireTime-time.Now().Unix())*time.Second)

	return nil

}

func GetRedisSessionKey(cookie string) string {

	return strings.Replace(conf.Get("redis_session_key").(string), "{cookie}", cookie, 1)
}
