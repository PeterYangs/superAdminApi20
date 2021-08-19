package task

import (
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

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}

func (t *BaseTask) Bind(taskParams interface{}, p map[string]interface{}) {

	s := reflect.ValueOf(taskParams).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
