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

			now := time.Now()

			if s.month != nil {

				continue
			}

			if s.minute != nil {

				if now.Second() == 0 {

					if s.minute.every {

						if now.Minute()%s.minute.value == 0 {

							go s.fn()
						}

					} else {

						if now.Minute() == s.minute.value {

							go s.fn()
						}

					}

				}

				continue
			}

		}

		fmt.Println("---------------")

		time.Sleep(1 * time.Second)

	}

}
