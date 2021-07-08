package logs

import (
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/spf13/cast"
	"log"
	"os"
	"runtime"
	"runtime/debug"
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
	queue     chan logs
	logLevels map[level]*logLevel
}

type logLevel struct {
	filePath string
	file     *os.File
}

type logs struct {
	level   level
	message string
}

var service logsService

func CreateLogs() *logsService {

	service = logsService{
		queue: make(chan logs, 10),
		logLevels: map[level]*logLevel{
			Error: {
				filePath: "logs/error.log",
			},
			Info: {
				filePath: "logs/info.log",
			},
			Debug: {
				filePath: "logs/debug.log",
			},
		},
	}

	return &service

}

func NewLogs() *logsService {

	return &service
}

func (l *logsService) Error(message string) *result {

	m := logFormat(Error, message)

	l.queue <- logs{
		level:   Error,
		message: m,
	}

	return &result{
		message: m,
	}
}

func (l *logsService) Info(message string) *result {

	m := logFormat(Info, message)

	l.queue <- logs{
		level:   Info,
		message: m,
	}

	return &result{
		message: m,
	}
}

func (l *logsService) Debug(message string) *result {

	m := logFormat(Debug, message)

	l.queue <- logs{
		level:   Debug,
		message: m,
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

func (l *logsService) Task() {

	defer close(l.queue)

	for _, item := range l.logLevels {

		f, e := os.OpenFile(item.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)

		if e != nil {

			log.Fatal(e)

		}

		item.file = f

	}

	for message := range l.queue {

		l.logLevels[message.level].file.Write([]byte(message.message))

	}

}

func logFormat(level level, message string) string {

	//f, l := GetLine()

	_, f, l, _ := runtime.Caller(2)

	m := tools.Date("[Y-m-d H:i:s ]", time.Now().Unix()) + " " + f + ":" + cast.ToString(l) + " " + level.ToString() + " : " + message + "\n"

	if level == Error {

		m += " \n[stacktrace]\n" + string(debug.Stack()) + "\n\n"
	}

	return m

}
