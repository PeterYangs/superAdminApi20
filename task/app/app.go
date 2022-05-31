package app

import (
	"github.com/PeterYangs/superAdminCore/v2/queue/task"
)

type AppTask struct {
	task.BaseTask
	Parameters *Parameter
}

type Parameter struct {
	task.Parameter
}

func NewAppTask() *AppTask {

	return &AppTask{

		BaseTask: task.BaseTask{
			TaskName: "app",
		},
		Parameters: &Parameter{},
	}
}

func (t *AppTask) Run() error {

	return nil
}

func (t *AppTask) BindParameters(p map[string]interface{}) {

	t.BaseTask.Bind(t.Parameters, p)

}
