package logs

import (
	"context"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/spf13/cast"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type level string

func (l level) ToString() string {

	return string(l)
}

const (
	Error level = "Error"
	Info  level = "Info"
	Debug level = "Debug"
)

type logsService struct {
	queue     chan *logs
	logLevels map[level]*logLevel
}

type logLevel struct {
	fileDir string
	file    *os.File
}

type logs struct {
	level   level
	message string
	time    time.Time
	lock    sync.Mutex
}

var service *logsService

func CreateLogs() *logsService {

	service = &logsService{
		queue: make(chan *logs, 10),
		logLevels: map[level]*logLevel{
			Error: {
				fileDir: "logs/error",
			},
			Info: {
				fileDir: "logs/info",
			},
			Debug: {
				fileDir: "logs/debug",
			},
		},
	}

	service.MakeDir()

	return service

}

func (ls *logsService) MakeDir() {

	for _, l2 := range ls.logLevels {

		os.MkdirAll(l2.fileDir, 0755)

	}
}

func NewLogs() *logsService {

	return service
}

func (ls *logsService) Error(message string) *result {

	m := logFormat(Error, message)

	ls.queue <- &logs{
		level:   Error,
		message: m,
		time:    time.Now(),
	}

	return &result{
		message: m,
	}
}

func (ls *logsService) Info(message string) *result {

	m := logFormat(Info, message)

	ls.queue <- &logs{
		level:   Info,
		message: m,
		time:    time.Now(),
	}

	return &result{
		message: m,
	}
}

func (ls *logsService) Debug(message string) *result {

	m := logFormat(Debug, message)

	ls.queue <- &logs{
		level:   Debug,
		message: m,
		time:    time.Now(),
	}

	return &result{
		message: m,
	}
}

type result struct {
	message string
}

func (r *result) Message() string {

	return r.message
}

func (r *result) Stdout() {

	fmt.Println(r.message)
}

func (ls *logsService) Task(cxt context.Context, wait *sync.WaitGroup) {

	defer wait.Done()

	go func() {

		select {

		case <-cxt.Done():

			close(ls.queue)

		}

	}()

	for message := range ls.queue {

		if ls.logLevels[message.level].fileDir == "" {

			continue
		}

		message.checkFileName(ls.logLevels[message.level])

		ls.logLevels[message.level].file.Write([]byte(message.message))

	}

	fmt.Println("日志模块安全退出")

}

func logFormat(level level, message string) string {

	_, f, l, _ := runtime.Caller(2)

	m := tools.Date("[Y-m-d H:i:s ]", time.Now().Unix()) + " " + f + ":" + cast.ToString(l) + " " + level.ToString() + " : " + message + "\n"

	if level == Error {

		m += " \n[stacktrace]\n" + string(debug.Stack()) + "\n\n"
	}

	return m

}

func (logs *logs) checkFileName(logLevel *logLevel) {

	logs.lock.Lock()

	defer logs.lock.Unlock()

	name := logs.fileFormat()

	if logLevel.file == nil || logLevel.file.Name() != name {

		//关闭文件(由于文件是根据日期来着，过了一天要关闭上一天的文件流)
		logLevel.file.Close()

		f, e := os.OpenFile(logLevel.fileDir+"/"+name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

		if e != nil {

			return

			//panic(e)

		}

		logLevel.file = f

	}

}

func (logs *logs) fileFormat() string {

	return tools.Date("Y-m-d", logs.time.Unix()) + ".log"
}
