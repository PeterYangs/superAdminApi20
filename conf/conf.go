package conf

import (
	"os"
	"sync"
)

var _conf map[string]interface{}

var lock sync.Mutex

func Load() {

	_conf = map[string]interface{}{

		"cookie_name":       os.Getenv("APP_NAME") + "_session", //浏览器cookie名称
		"cookie_key":        "cookie_key",                       //context中cookie的值的name
		"redis_prefix":      os.Getenv("APP_NAME") + ":",        //redis前缀
		"redis_session_key": "session:{cookie}",                 //session在redis中的key(带redis前缀)
		"captcha_key":       "_captcha",                         //验证码的key
		"captcha_lifetime":  os.Getenv("CAPTCHA_LIFETIME"),      //验证码过期时间
		"lock_prefix":       os.Getenv("LOCK_PREFIX"),           //锁前缀
	}

}

func Get(key string) interface{} {

	lock.Lock()

	defer lock.Unlock()

	return _conf[key]

}
