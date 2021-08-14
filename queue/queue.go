package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-web/component/logs"
	"gin-web/interface/task"
	"gin-web/redis"
	"gin-web/task/email"
	"gin-web/task/sms"
	"github.com/PeterYangs/tools"
	redis2 "github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

var handles = sync.Map{}

type job struct {
	Delay_ time.Duration `json:"-"` //延迟
	Data_  task.Task     `json:"data"`
	Queue_ string        `json:"queue"` //队列名称
	Id     string        `json:"id"`
}

var once sync.Once

func init() {

	handles.Store("email", &email.TaskEmail{Parameters: &email.Parameter{}})
	handles.Store("sms", &sms.TaskSms{Parameters: &sms.Parameter{}})

}

func Run() {

	defer func() {
		if r := recover(); r != nil {

			fmt.Println(r)

			fmt.Println(string(debug.Stack()))

			msg := fmt.Sprint(r)

			logs.NewLogs().Error(msg)

		}
	}()

	once.Do(func() {

		go checkDelay()
	})

	queues := tools.Explode(",", os.Getenv("QUEUES"))

	for i, queue := range queues {

		queues[i] = os.Getenv("QUEUE_PREFIX") + queue
	}

	for {

		//timeout为0则为永久超时
		s, err := redis.GetClient().BLPop(context.TODO(), 0, queues...).Result()

		if err != nil {

			log.Println(err)

			fmt.Println("队列退出")

			break
		}

		var jsons map[string]interface{}

		err = json.Unmarshal([]byte(s[1]), &jsons)

		if err != nil {

			fmt.Println(err)

			continue
		}

		data := jsons["data"].(map[string]interface{})

		//fmt.Println(p)

		////获取task
		hh, ok := handles.Load(data["TaskName"].(string))
		//hh, ok := handles.Load(jsons.Data_.GetName())
		//
		h := hh.(task.Task)

		if !ok {

			fmt.Println("获取task失败")

			continue
		}

		////绑定参数
		h.BindParameters(cast.ToStringMapString(data["Parameters"]))
		//
		////执行任务
		h.Run()

	}

}

func checkDelay() {

	defer func() {
		if r := recover(); r != nil {

			fmt.Println(r)

			fmt.Println(string(debug.Stack()))

			msg := fmt.Sprint(r)

			logs.NewLogs().Error(msg)

		}
	}()

	for {

		push()

		time.Sleep(1 * time.Second)

	}

}

func push() {

	//分布式锁
	lock := redis.GetClient().Lock("queue:delay:lock", 10*time.Second)

	defer lock.Release()

	if !lock.Get() {

		time.Sleep(1 * time.Second)

		return
	}

	//查询已到期任务
	list, err := redis.GetClient().ZRangeByScore(context.TODO(), os.Getenv("QUEUE_PREFIX")+"delay", &redis2.ZRangeBy{
		Min: "0",
		Max: cast.ToString(time.Now().Unix()),
	}).Result()

	if err != nil {

		fmt.Println(err)

		time.Sleep(1 * time.Second)
		return
	}

	for _, s := range list {

		var jsons map[string]interface{}

		json.Unmarshal([]byte(s), &jsons)

		queue := ""

		if jsons["queue"].(string) == "" {

			queue = os.Getenv("QUEUE_PREFIX") + os.Getenv("DEFAULT_QUEUE")

		} else {

			queue = os.Getenv("QUEUE_PREFIX") + jsons["queue"].(string)
		}

		//头部插入，先执行
		redis.GetClient().LPush(context.TODO(), queue, s).Result()

	}

	if len(list) > 0 {

		//删除已到期的任务
		redis.GetClient().ZRemRangeByRank(context.TODO(), os.Getenv("QUEUE_PREFIX")+"delay", 0, int64(len(list)-1))
	}

}

func Dispatch(task task.Task) *job {

	//t, _ := json.Marshal(task)

	return &job{
		Data_:  task,
		Delay_: 0,
		Id:     uuid.NewV4().String(),
	}

}

func (j *job) Delay(duration time.Duration) *job {

	j.Delay_ = duration

	return j
}

func (j *job) Queue(queue string) *job {

	j.Queue_ = queue

	return j
}

func (j *job) Run() {

	queue := ""

	if j.Delay_ == 0 {

		if j.Queue_ == "" {

			queue = os.Getenv("QUEUE_PREFIX") + os.Getenv("DEFAULT_QUEUE")

		} else {

			queue = os.Getenv("QUEUE_PREFIX") + j.Queue_
		}

		data, err := json.Marshal(j)

		if err != nil {

			return
		}

		redis.GetClient().RPush(context.TODO(), queue, data)

	} else {

		if j.Queue_ == "" {

			queue = os.Getenv("QUEUE_PREFIX") + "delay"

		} else {

			queue = os.Getenv("QUEUE_PREFIX") + "delay"
		}

		//queueName := os.Getenv("DEFAULT_QUEUE")
		//
		//if j.queue != "" {
		//
		//	queueName = j.queue
		//}

		//json.Marshal()

		data, err := json.Marshal(j)

		//fmt.Println()
		//
		//fmt.Println(string(data))

		if err != nil {

			fmt.Println(err)

			return
		}

		redis.GetClient().ZAdd(context.TODO(), queue, &redis2.Z{
			Score:  float64(time.Now().Unix() + cast.ToInt64(j.Delay_.Seconds())),
			Member: data,
		})
	}

}
