package task

type Task interface {
	Run()
	BindParameters(map[string]string)
	GetName() string
}
