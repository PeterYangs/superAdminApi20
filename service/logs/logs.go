package logs

import error2 "gin-web/service/logs/error"

type Level string

type Logs interface {
	Log(message string)
	GetFilePath() string
	GetLevel() Level
	GetMessage() string
}

type logsService struct {
	Queue  chan Logs
	levels []Logs
}

var Service *logsService

func CreateLogs() *logsService {

	//e:=

	Service = &logsService{
		Queue: make(chan Logs, 10),
		levels: []Logs{

			error2.Registered(""),
		},
	}

	return Service

}
