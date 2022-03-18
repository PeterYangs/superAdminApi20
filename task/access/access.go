package access

import (
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/queue/task"
	"superadmin/model"
)

type TaskAccess struct {
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

func NewAccessTask(ip string, url string, params string, adminId float64) *TaskAccess {

	return &TaskAccess{

		BaseTask: task.BaseTask{
			TaskName: "access",
		},
		Parameters: &Parameter{
			Ip:      ip,
			Url:     url,
			Params:  params,
			AdminId: adminId,
		},
	}
}

func (t *TaskAccess) Run() error {

	database.GetDb().Create(&model.Access{
		Ip:      t.Parameters.Ip,
		Url:     t.Parameters.Url,
		Params:  t.Parameters.Params,
		AdminId: t.Parameters.AdminId,
	})

	return nil
}

func (t *TaskAccess) BindParameters(p map[string]interface{}) {

	t.BaseTask.Bind(t.Parameters, p)

}
