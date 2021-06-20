package main

import (
	"context"
	"fmt"
	"gin-web/redis"
)

func main() {

	//redis.Connect()

	//fmt.Println(redis.Client)

	redis.GetClient().Set(context.TODO(), "hh", "yy", 0)

	re, _ := redis.GetClient().Get(context.TODO(), "hh").Result()

	fmt.Println(re)
}
