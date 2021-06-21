package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
	"time"
)

//var client *redis.Client

var once = sync.Once{}

type _connect struct {
	connect *redis.Client
}

var Client *_connect

func connect() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	conf := &redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   0,
	}

	if os.Getenv("REDIS_PASSWORD") != "null" {

		conf.Password = os.Getenv("REDIS_PASSWORD")
	}

	//c := redis.NewClient(conf)

	c := redis.NewClient(conf)

	re := c.Ping(cxt)

	defer cancel()

	if re.Err() != nil {

		//cancel()

		panic(re.Err())

	}

	//client = c

	//cc.Connect=c

	Client = &_connect{
		connect: c,
	}

}

func GetClient() *_connect {

	//fmt.Println(Client)

	if Client == nil {

		once.Do(func() {

			connect()
		})

	}

	return Client

}

func (cc *_connect) Get(cxt context.Context, key string) *redis.StringCmd {

	return cc.connect.Get(cxt, key)
}

func (cc *_connect) Set(cxt context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {

	return cc.connect.Set(cxt, key, value, expiration)
}
