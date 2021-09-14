package queue

import (
	"github.com/PeterYangs/superAdminCore/queue/template"
	"superadmin/task/access"
)

var Queues = map[string]template.Task{
	"access": &access.TaskAccess{Parameters: &access.Parameter{}},
}
