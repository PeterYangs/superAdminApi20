package task

import (
	"fmt"
	"reflect"
)

type BaseTask struct {
	TaskName   string //处理器名称
	Parameters *Parameter
}

type Parameter struct {
}

func (t *BaseTask) Run() {

}

func (t *BaseTask) GetName() string {

	return t.TaskName
}

func (t *BaseTask) BindParameters(p map[string]interface{}) {

	//t.Parameters= p

	fmt.Println(t)

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
