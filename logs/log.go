package logs

import (
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/spf13/cast"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
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

type logs struct {
	queue chan string
	file  *os.File
}

var Logs logs

func CreateLogs(filePath string) *logs {

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)

	if err != nil {

		panic(err)

	}

	Logs = logs{
		queue: make(chan string, 10),
		file:  f,
	}

	return &Logs

}

func NewLogs() *logs {

	return &Logs
}

func (l *logs) Error(message string) *result {

	m := logFormat(Error, message)

	l.queue <- m

	return &result{
		message: m,
	}
}

func (l *logs) Info(message string) *result {

	m := logFormat(Info, message)

	l.queue <- m

	return &result{
		message: m,
	}
}

func (l *logs) Debug(message string) *result {

	m := logFormat(Debug, message)

	l.queue <- m

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

func (l *logs) Task() {

	defer func() {

		close(l.queue)

		l.file.Close()

	}()

	for message := range l.queue {

		l.file.Write([]byte(message))
	}

}

func logFormat(level level, message string) string {

	f, l := GetLine()

	m := tools.Date("[Y-m-d H:i:s ]", time.Now().Unix()) + " " + f + ":" + cast.ToString(l) + " " + level.ToString() + " : " + message

	if level == Error {

		m += " \n[stacktrace]\n" + string(debug.Stack()) + "\n\n"
	}

	return m

}

func GetLine() (string, int) {

	for i := 0; true; i++ {

		pc, _, _, ok := runtime.Caller(i)

		if ok {

			if strings.Contains(runtime.FuncForPC(pc).Name(), "logs.(*logs)") {

				ii := i + 1

				_, f, l, _ := runtime.Caller(ii)

				return f, l

			}

		} else {

			return "", 0
		}

	}

	return "", 0
}
