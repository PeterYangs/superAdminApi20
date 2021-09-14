package queue

import (
	"github.com/PeterYangs/superAdminCore/queue/template"
	"superadmin/task/access"
	"superadmin/task/app"
	//namespace
)

var Queues = map[string]template.Task{
	"access": &access.TaskAccess{Parameters: &access.Parameter{}},
	"app":    &app.AppTask{Parameters: &app.Parameter{}},
	//taskRegister
}
