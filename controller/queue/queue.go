package queue

import (
	"context"
	"encoding/json"
	"gin-web/contextPlus"
	"gin-web/redis"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"os"
)

func List(c *contextPlus.Context) *response.Response {

	p := cast.ToInt(c.DefaultQuery("p", "1"))

	queue := c.DefaultQuery("queue", os.Getenv("DEFAULT_QUEUE"))

	size := 10

	start := (p - 1) * size

	stop := start + size

	list, _ := redis.GetClient().LRange(context.TODO(), os.Getenv("QUEUE_PREFIX")+queue, int64(start), int64(stop)).Result()

	count, _ := redis.GetClient().LLen(context.TODO(), os.Getenv("QUEUE_PREFIX")+queue).Result()

	if len(list) <= 0 {

		return response.Resp().Api(1, "success", gin.H{"total": 0, "data": []string{}, "page": 1, "size": size})
	}

	jList := make([]map[string]interface{}, len(list))

	for i, s := range list {

		temp := make(map[string]interface{})

		json.Unmarshal([]byte(s), &temp)

		jList[i] = temp
	}

	return response.Resp().Api(1, "success", gin.H{"total": count, "data": jList, "page": p, "size": size})

}

func DelayList(c *contextPlus.Context) *response.Response {

	p := cast.ToInt(c.DefaultQuery("p", "1"))

	//queue := c.DefaultQuery("queue", os.Getenv("DEFAULT_QUEUE"))

	size := 10

	start := (p - 1) * size

	stop := start + size

	list, _ := redis.GetClient().ZRange(context.TODO(), os.Getenv("QUEUE_PREFIX")+"delay", int64(start), int64(stop)).Result()

	count, _ := redis.GetClient().ZCard(context.TODO(), os.Getenv("QUEUE_PREFIX")+"delay").Result()

	if len(list) <= 0 {

		return response.Resp().Api(1, "success", gin.H{"total": 0, "data": []string{}, "page": 1, "size": size})
	}

	jList := make([]map[string]interface{}, len(list))

	for i, s := range list {

		temp := make(map[string]interface{})

		json.Unmarshal([]byte(s), &temp)

		jList[i] = temp
	}

	return response.Resp().Api(1, "success", gin.H{"total": count, "data": jList, "page": p, "size": size})

}
