package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var client *redis.Client

var once = sync.Once{}

func Connect() {

	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	c := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

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

			Connect()
		})

	}

	return client

}
