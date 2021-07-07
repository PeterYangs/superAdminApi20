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

type logsService struct {
	queue chan string
	//file    *os.File
	logPath string
	errPath string
}

var service logsService

func CreateLogs(logPath string, errPath string) *logsService {

	//f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)
	//
	//if err != nil {
	//
	//	panic(err)
	//
	//}

	service = logsService{
		queue:   make(chan string, 10),
		logPath: logPath,
		errPath: errPath,
	}

	return &service

}

func NewLogs() *logsService {

	return &service
}

func (l *logsService) Error(message string) *result {

	m := logFormat(Error, message)

	l.queue <- m

	return &result{
		message: m,
	}
}

func (l *logsService) Info(message string) *result {

	m := logFormat(Info, message)

	l.queue <- m

	return &result{
		message: m,
	}
}

func (l *logsService) Debug(message string) *result {

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

func (l *logsService) Task() {

	fileIsOpen := true

	defer close(l.queue)

	logfile, err := os.OpenFile(l.logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)

	if err != nil {

		fileIsOpen = false

		fmt.Println(err)
	}

	defer logfile.Close()

	errfile, err := os.OpenFile(l.errPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)

	if err != nil {

		fileIsOpen = false

		fmt.Println(err)
	}

	defer errfile.Close()

	for message := range l.queue {

		if !fileIsOpen {

			fmt.Println(message)

		}

		//l.file.Write([]byte(message))
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
