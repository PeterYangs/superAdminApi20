package register

import (
	"gin-web/task/access"
	"gin-web/task/email"
	"gin-web/task/sms"
	//namespace
	"sync"
)

type handles struct {
	Tasks map[string]interface{}
	Lock  sync.Mutex
}

var Register = map[string]interface{}{
	"email":  &email.TaskEmail{Parameters: &email.Parameter{}},
	"sms":    &sms.TaskSms{Parameters: &sms.Parameter{}},
	"access": &access.TaskAccess{Parameters: &access.Parameter{}},
	//taskRegister
}

var Handles = &handles{
	Tasks: make(map[string]interface{}),
	Lock:  sync.Mutex{},
}

func (h *handles) Init() {

	h.Lock.Lock()

	defer h.Lock.Unlock()

	h.Tasks = Register

}

func (h *handles) GetTask(name string) (interface{}, bool) {

	h.Lock.Lock()

	defer h.Lock.Unlock()

	task, ok := h.Tasks[name]

	return task, ok
}
