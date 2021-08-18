package access

import (
	"gin-web/database"
	"gin-web/model"
	"reflect"
)

type TaskAccess struct {
	TaskName   string //处理器名称
	Parameters *Parameter
}

type Parameter struct {
	Ip     string
	Url    string
	Params string
}

func NewTask(ip string, url string, params string) *TaskAccess {

	return &TaskAccess{
		Parameters: &Parameter{
			Ip:     ip,
			Url:    url,
			Params: params,
		},
		TaskName: "access",
	}
}

func (t *TaskAccess) Run() {

	//time.Sleep(200*time.Millisecond)
	//fmt.Println(t.Parameters.Email, t.Parameters.Title, t.Parameters.Content)

	//fmt.Println(1111)

	database.GetDb().Create(&model.Access{
		Ip:     t.Parameters.Ip,
		Url:    t.Parameters.Url,
		Params: t.Parameters.Params,
	})

}

func (t *TaskAccess) GetName() string {

	return t.TaskName
}

func (t *TaskAccess) BindParameters(p map[string]string) {

	//t.Parameters= p

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
