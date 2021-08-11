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
	"github.com/spf13/cast"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var handles = sync.Map{}

type job struct {
	delay time.Duration //延迟
	data  []byte
}

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

	for {

		s, err := redis.GetClient().BRPop(context.TODO(), 0, "queue:default").Result()

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

		////获取task
		hh, ok := handles.Load(jsons["TaskName"].(string))

		h := hh.(task.Task)

		if !ok {

			fmt.Println("获取task失败")

			continue
		}

		//绑定参数
		h.BindParameters(cast.ToStringMapString(jsons["Parameters"]))

		//执行任务
		h.Run()

	}

	//eval
}

func Dispatch(task task.Task) *job {

	t, _ := json.Marshal(task)

	return &job{
		data:  t,
		delay: 0,
	}

}

func (j *job) Delay(duration time.Duration) *job {

	j.delay = duration

	return j
}

func (j *job) Run() {

	//fmt.Println(j.delay.Seconds())

	redis.GetClient().LPush(context.TODO(), "queue:default", j.data)
}
