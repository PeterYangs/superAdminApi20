package task

type Task interface {
	Run()
	BindParameters(map[string]interface{})
	GetName() string
}
