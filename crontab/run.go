package crontab

import (
	"time"
)

func Run() {

	_crontab := &crontab{}

	Registered(_crontab)

	start := true

	delay := true

	for {

		//每小时校对一下时间
		if delay || time.Now().Minute() == 0 {

			for {

				if time.Now().Second() == 0 {

					//fmt.Println("do")

					delay = false

					break
				}

				//fmt.Println(now.Second())

				time.Sleep(1 * time.Second)

			}

		}

		if !start {

			time.Sleep(1 * time.Minute)
		}

		now := time.Now()

		for _, s := range _crontab.schedules {

			if s.day != nil {

				dealDay(s, now)

				continue
			}

			if s.hour != nil {

				dealHour(s, now)

				continue
			}

			if s.minute != nil {

				dealMinute(s, now)

				continue
			}

		}

		start = false

	}

}

func dealMinute(s *schedule, now time.Time) {

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

func dealHour(s *schedule, now time.Time) {

	if s.minute == nil {

		if now.Minute() == 0 {

			if s.hour.every {

				if now.Hour()%s.hour.value == 0 {

					go s.fn()
				}

			} else {

				if now.Hour() == s.hour.value {

					go s.fn()
				}

			}

		}

	} else {

		if s.hour.every {

			if now.Hour()%s.hour.value == 0 {

				//go s.fn()

				dealMinute(s, now)

			}

		} else {

			if now.Hour() == s.hour.value {

				//go s.fn()

				dealMinute(s, now)
			}

		}

	}

}

func dealDay(s *schedule, now time.Time) {

	if s.hour == nil {

		if now.Hour() == 0 {

			if s.hour.every {

				if now.Day()%s.day.value == 0 {

					go s.fn()
				}

			} else {

				if now.Day() == s.day.value {

					go s.fn()
				}

			}

		}

	} else {

		if s.day.every {

			if now.Day()%s.day.value == 0 {

				//go s.fn()

				dealHour(s, now)

			}

		} else {

			if now.Day() == s.day.value {

				//go s.fn()

				dealHour(s, now)
			}

		}

	}

}
