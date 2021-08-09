package main

import (
	"fmt"
	"gin-web/migrate/bin"
	bin2 "gin-web/seeds/bin"
	"github.com/manifoldco/promptui"
)

func main() {

	prompt := promptui.Select{
		Label: "选择类型",
		Items: []string{"数据库迁移", "数据填充"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {

	case "数据库迁移":

		new(bin.MigrateRun).Run()

	case "数据填充":

		new(bin2.Bin).Run()

	}
}
