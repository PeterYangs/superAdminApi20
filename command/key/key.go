package key

import (
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PeterYangs/tools/file/read"
	uuid "github.com/satori/go.uuid"
	"os"
	"regexp"
)

type Key struct {
}

func (k Key) Run() {

	_, err := os.Stat(".env")

	if err != nil {
		panic(err)
	}

	if os.IsNotExist(err) {

		fmt.Println(".env文件不存在")

		return
	}

	res, err := read.Open(".env").Read()

	if err != nil {

		panic(err)
	}

	re1 := regexp.MustCompile("KEY=[0-9A-Za-z!@#$%^&*]+").ReplaceAllString(string(res), "KEY="+tools.Md5(uuid.NewV4().String()))

	f, err := os.OpenFile(".env", os.O_RDWR, 0644)

	if err != nil {

		panic(err)
	}

	defer f.Close()

	_, err = f.Write([]byte(re1))

	if err != nil {

		panic(err)
	}

}
