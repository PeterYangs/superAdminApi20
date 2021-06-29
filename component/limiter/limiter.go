package limiter

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type KeyLimiter struct {
	limiter *rate.Limiter
	lastGet time.Time //上一次获取token的时间
	key     string
}

type Limiters struct {
	limiter map[string]*KeyLimiter
	lock    sync.Mutex
}

var GlobalLimiters = &Limiters{
	limiter: make(map[string]*KeyLimiter),
	lock:    sync.Mutex{},
}

func NewLimiter(r rate.Limit, b int, key string) *KeyLimiter {

	go GlobalLimiters.clearLimiter()

	keyLimiter := GlobalLimiters.Get(r, b, key)

	return keyLimiter

}

func (l *KeyLimiter) Allow() bool {

	//fmt.Println(globalLimiters)

	//fmt.Println(l.lastGet)

	l.lastGet = time.Now()

	return l.limiter.Allow()

}

func (ls *Limiters) Get(r rate.Limit, b int, key string) *KeyLimiter {

	ls.lock.Lock()

	defer ls.lock.Unlock()

	limiter, ok := ls.limiter[key]

	if ok {

		return limiter
	}

	l := &KeyLimiter{
		limiter: rate.NewLimiter(r, b),
		lastGet: time.Now(),
		key:     key,
	}

	ls.limiter[key] = l

	return l
}

//清除过期的限流器
func (ls *Limiters) clearLimiter() {

	for {

		time.Sleep(1 * time.Minute)

		for i, i2 := range ls.limiter {

			//超过1分钟
			if time.Now().Unix()-i2.lastGet.Unix() > 60 {
				ls.lock.Lock()
				delete(ls.limiter, i)
				ls.lock.Unlock()
			}

		}

	}

}
