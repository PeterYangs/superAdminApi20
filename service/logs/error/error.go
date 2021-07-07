package error

import "gin-web/service/logs"

const level logs.Level = "Error"

type Error struct {
	filePath string
	level    logs.Level
	message  string
}

func Registered(filePath string) *Error {

	return &Error{filePath: filePath, level: level}

}

func New() *Error {

	return &Error{level: level}
}

func (e *Error) Log(message string) {

	e.message = message

	logs.Service.Queue <- e

}

func (e *Error) GetFilePath() string {

	return e.filePath
}

func (e *Error) GetLevel() logs.Level {

	return e.level
}

func (e *Error) GetMessage() string {

	return e.message
}
