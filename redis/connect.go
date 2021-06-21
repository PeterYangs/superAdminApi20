package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
	"time"
)

var client *redis.Client

var once = sync.Once{}

func connect() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	conf := &redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   0,
	}

	if os.Getenv("REDIS_PASSWORD") != "null" {

		conf.Password = os.Getenv("REDIS_PASSWORD")
	}

	c := redis.NewClient(conf)

	re := c.Ping(cxt)

	defer cancel()

	if re.Err() != nil {

		//cancel()

		panic(re.Err())

	}

	client = c

}

func GetClient() *redis.Client {

	if client == nil {

		once.Do(func() {

			connect()
		})

	}

	return client

}
