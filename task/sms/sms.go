package sms

import (
	"fmt"
	"reflect"
)

type TaskSms struct {
	TaskName   string //处理器名称
	Parameters *Parameter
}

type Parameter struct {
	Phone   string
	Content string
}

func NewTask(phone, content string) *TaskSms {

	return &TaskSms{
		Parameters: &Parameter{
			Phone:   phone,
			Content: content,
		},
		TaskName: "sms",
	}
}

func (t *TaskSms) Run() {

	fmt.Println(t.Parameters.Phone, t.Parameters.Content)

}

func (t *TaskSms) GetName() string {

	return t.TaskName
}

func (t *TaskSms) BindParameters(p map[string]string) {

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
