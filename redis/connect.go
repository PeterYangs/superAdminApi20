package redis

import (
	"context"
	"gin-web/conf"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"os"
	"sync"
	"time"
)

//var client *redis.Client

var once = sync.Once{}

type _connect struct {
	connect *redis.Client
}

type lock struct {
	key         string
	expiration  time.Duration
	connect     *_connect
	requestId   string
	checkCancel chan bool
	mu          sync.Mutex
}

var Client *_connect

func connect() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	conf_ := &redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   0,
	}

	if os.Getenv("REDIS_PASSWORD") != "null" {

		conf_.Password = os.Getenv("REDIS_PASSWORD")
	}

	c := redis.NewClient(conf_)

	re := c.Ping(cxt)

	defer cancel()

	if re.Err() != nil {

		panic(re.Err())

	}

	Client = &_connect{
		connect: c,
	}

}

func GetClient() *_connect {

	if Client == nil {

		once.Do(func() {

			connect()
		})

	}

	return Client

}

func (cc *_connect) Get(cxt context.Context, key string) *redis.StringCmd {

	return cc.connect.Get(cxt, conf.Get("redis_prefix").(string)+key)
}

func (cc *_connect) Set(cxt context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {

	return cc.connect.Set(cxt, conf.Get("redis_prefix").(string)+key, value, expiration)
}

func (cc *_connect) Exists(cxt context.Context, keys ...string) *redis.IntCmd {

	//cc.connect.Expire()

	for i, key := range keys {

		keys[i] = conf.Get("redis_prefix").(string) + key
	}

	return cc.connect.Exists(cxt, keys...)
}

func (cc *_connect) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {

	return cc.connect.Expire(ctx, conf.Get("redis_prefix").(string)+key, expiration)
}

func (cc *_connect) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {

	return cc.connect.SetNX(ctx, conf.Get("redis_prefix").(string)+key, value, expiration)

}

func (cc *_connect) Del(ctx context.Context, keys ...string) *redis.IntCmd {

	for i, key := range keys {

		keys[i] = conf.Get("redis_prefix").(string) + key
	}

	return cc.connect.Del(ctx, keys...)
}

func (cc *_connect) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {

	return cc.connect.Eval(ctx, script, keys, args...)
}

// Lock 声明锁
func (cc *_connect) Lock(key string, expiration time.Duration) *lock {

	lockId := uuid.NewV4().String()

	return &lock{key: conf.Get("lock_prefix").(string) + key, expiration: expiration, connect: cc, requestId: lockId, mu: sync.Mutex{}}
}

// Get 获取锁
func (lk *lock) Get() bool {

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	ok, err := lk.connect.SetNX(cxt, lk.key, lk.requestId, lk.expiration).Result()

	if err != nil {

		return false
	}

	if ok {

		//检查锁是否过期
		go lk.checkLockIsRelease()

	}

	return ok
}

// Block 获取锁阻塞
func (lk *lock) Block(expiration time.Duration) bool {

	t := time.Now()

	for {

		cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		ok, err := lk.connect.SetNX(cxt, lk.key, lk.requestId, lk.expiration).Result()

		cancel()

		if err != nil {

			return false
		}

		if ok {

			go lk.checkLockIsRelease()

			return true
		}

		time.Sleep(200 * time.Millisecond)

		if time.Now().Sub(t) > expiration {

			return false
		}

	}

}

// Release 释放锁
func (lk *lock) Release() error {

	lk.mu.Lock()

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	const luaScript = `
	if redis.call('get', KEYS[1])==ARGV[1] then
		return redis.call('del', KEYS[1])
	else
		return 0
	end
	`
	res, err := redis.NewScript(luaScript).Run(cxt, GetClient().connect, []string{conf.Get("redis_prefix").(string) + lk.key}, lk.requestId).Result()

	lk.mu.Unlock()

	if err != nil {

		return err
	}

	if res.(int64) != 0 {

		lk.checkCancel <- true
	}

	return err

}

// ForceRelease 强制释放锁，忽略请求id
func (lk *lock) ForceRelease() error {

	lk.mu.Lock()

	defer func() {

		lk.mu.Unlock()

		if lk.checkCancel != nil {

			lk.checkCancel <- true
		}

	}()

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := lk.connect.Del(cxt, lk.key).Result()

	return err

}

//检查锁是否被释放，未被释放就延长锁时间
func (lk *lock) checkLockIsRelease() {

	for {

		checkCxt, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(lk.expiration.Milliseconds()-lk.expiration.Milliseconds()/10))

		lk.checkCancel = make(chan bool)

		select {

		case <-checkCxt.Done():

			//多次续期，直到锁被释放
			isContinue := lk.done()

			if !isContinue {

				return
			}

		//取消
		case <-lk.checkCancel:

			//fmt.Println("释放")

			return

		}

	}
}

//判断锁是否已被释放
func (lk *lock) done() bool {

	lk.mu.Lock()

	defer lk.mu.Unlock()

	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	res, err := lk.connect.Exists(cxt, lk.key).Result()

	cancel()

	if err != nil {

		return false
	}

	if res == 1 {

		cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		ok, err := lk.connect.Expire(cxt, lk.key, lk.expiration).Result()

		cancel()

		if err != nil {

			return false
		}

		if ok {

			//fmt.Println("续期")

			return true

		}

	}

	return false
}
