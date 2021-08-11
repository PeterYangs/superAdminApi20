package email

import (
	"fmt"
	"reflect"
)

//const Name taskName.TaskName = "email"

type TaskEmail struct {
	TaskName   string //处理器名称
	Parameters *Parameter
}

type Parameter struct {
	Email   string
	Title   string
	Content string
}

func NewTask(email string, title string, content string) *TaskEmail {

	return &TaskEmail{
		Parameters: &Parameter{
			Email:   email,
			Title:   title,
			Content: content,
		},
		TaskName: "email",
	}
}

func (t *TaskEmail) Run() {

	fmt.Println(t.Parameters.Email, t.Parameters.Title, t.Parameters.Content)

}

func (t *TaskEmail) GetName() string {

	return t.TaskName
}

func (t *TaskEmail) BindParameters(p map[string]string) {

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
