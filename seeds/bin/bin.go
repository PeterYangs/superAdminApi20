package bin

import (
	"fmt"
	"github.com/PeterYangs/gcmd"
	"github.com/manifoldco/promptui"
)

type Bin struct {
}

func (b Bin) Run() {

	prompt := promptui.Select{
		Label: "选择类型",
		Items: []string{"重置管理员", "生成基础菜单"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {

	case "重置管理员":

		gcmd.Command("go run seeds/rootCreate.go").Start()

	case "生成基础菜单":

		gcmd.Command("go run seeds/menuCreate.go").Start()

	}

}
