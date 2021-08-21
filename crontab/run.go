package crontab

import (
	"fmt"
	"time"
)

func Run() {

	_crontab := &crontab{}

	Registered(_crontab)

	for {

		for _, s := range _crontab.schedules {

			fmt.Println(s)
		}

		fmt.Println("---------------")

		time.Sleep(1 * time.Second)

	}

}
