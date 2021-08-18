package access

import (
	"gin-web/database"
	"gin-web/model"
	"gin-web/task"
	"reflect"
)

type TaskAccess struct {
	//TaskName   string //处理器名称
	//Parameters *Parameter
	task.BaseTask
	Parameters *Parameter
}

type Parameter struct {
	task.Parameter
	Ip      string
	Url     string
	Params  string
	AdminId float64
}

func NewTask(ip string, url string, params string, adminId float64) *TaskAccess {

	return &TaskAccess{

		//BaseTask.T

		//task.BaseTask.
		BaseTask: task.BaseTask{
			TaskName: "access",
		},
		Parameters: &Parameter{
			Ip:      ip,
			Url:     url,
			Params:  params,
			AdminId: adminId,
		},
		//TaskName: "access",
	}
}

func (t *TaskAccess) Run() {

	database.GetDb().Create(&model.Access{
		Ip:      t.Parameters.Ip,
		Url:     t.Parameters.Url,
		Params:  t.Parameters.Params,
		AdminId: t.Parameters.AdminId,
	})

}

//func (t *TaskAccess) GetName() string {
//
//	return t.TaskName
//}

func (t *TaskAccess) BindParameters(p map[string]interface{}) {

	//t.BindParameters(p)

	//t.Parameters= p

	s := reflect.ValueOf(t.Parameters).Elem()

	for key, value := range p {

		if s.FieldByName(key).IsValid() {

			s.FieldByName(key).Set(reflect.ValueOf(value))

		}

	}

}
