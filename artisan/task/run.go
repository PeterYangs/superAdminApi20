package task

import (
	"fmt"
	"gin-web/common"
	"gin-web/queue/register"
	"github.com/PeterYangs/tools/file/read"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"regexp"
)

type TaskCmd struct {
}

func (t TaskCmd) Run() {

	prompt := promptui.Prompt{
		Label: "输入任务名",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	//fmt.Println(result)

	for name, _ := range register.Register {

		//fmt.Println(name,"---")

		if name == result {

			log.Println("该任务名已存在")

			return
		}

	}

	os.Mkdir("task/"+result, 0755)

	f, err := os.OpenFile("task/"+result+"/"+result+".go", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)

	if err != nil {

		panic(err)

	}

	defer f.Close()

	script := `
package ` + result + `
	
import (	

	"gin-web/task"

)
	
type Task` + common.Capitalize(result) + ` struct {
	task.BaseTask
	Parameters *Parameter
}
	
type Parameter struct {

	task.Parameter
	
}
	
func NewTask` + common.Capitalize(result) + `() *Task` + common.Capitalize(result) + ` {
	
	return &Task` + common.Capitalize(result) + `{
	
		BaseTask: task.BaseTask{
			TaskName: "` + result + `",
		},
			Parameters: &Parameter{
	
			},
		}
	}
	
func (t *Task` + common.Capitalize(result) + `) Run() {
	
	
	
}
	
func (t *Task` + common.Capitalize(result) + `) BindParameters(p map[string]interface{}) {
	
	t.BaseTask.Bind(t.Parameters, p)
	
}
	`

	_, err = f.Write([]byte(script))

	if err != nil {

		fmt.Println(err)
	}

	b, err := read.Open("queue/register/register.go").Read()

	//fmt.Println(string(b))

	//regexp.MustCompile()

	newScript := regexp.MustCompile("//namespace").ReplaceAllString(string(b), "\"gin-web/task/"+result+"\"\n"+"//namespace")
	newScript = regexp.MustCompile("//taskRegister").ReplaceAllString(newScript, `"`+result+`":   &`+result+`.Task`+common.Capitalize(result)+`{Parameters: &`+result+`.Parameter{}},`+"\n"+"//taskRegister")

	//fmt.Println(newScript)

	ff, err := os.OpenFile("queue/register/register.go", os.O_CREATE, 0644)

	if err != nil {

		//fmt.Println(err)

		log.Println(err)

		return
	}

	defer ff.Close()

	_, err = ff.Write([]byte(newScript))

	if err != nil {

		log.Println(err)

		return
	}

}
