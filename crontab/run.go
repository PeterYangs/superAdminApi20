package crontab

import (
	"sync"
	"time"
)

func Run(wait *sync.WaitGroup) {

	_crontab := &crontab{
		quitWait: wait,
	}

	Registered(_crontab)

	start := true

	delay := true

	var diff time.Duration

	for {

		//准点校对
		if delay {

			for {

				if time.Now().Second() == 0 {

					delay = false

					break
				}

				time.Sleep(1 * time.Second)

			}

		}

		if !start {

			//消除时间误差
			//time.Sleep(1*time.Minute)
			time.Sleep(1*time.Minute - diff)
		}

		startTime := time.Now()

		now := time.Now()

		//now:=time.Date(2021,9,12,0,1,0,0,time.Local)

		go deal(_crontab, now)

		start = false

		//计算时间误差
		diff = time.Now().Sub(startTime)

	}

}

func deal(crontab *crontab, now time.Time) {

	for _, s := range crontab.schedules {

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

				//时间区间
				if s.hour.between != nil {

					if now.Hour() >= s.hour.between.min && now.Hour() <= s.hour.between.max {

						go s.fn()
					}

				} else {

					if now.Hour() == s.hour.value {

						go s.fn()
					}
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

			//时间区间
			if s.hour.between != nil {

				if now.Hour() >= s.hour.between.min && now.Hour() <= s.hour.between.max {

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

}

func dealDay(s *schedule, now time.Time) {

	if s.hour == nil {

		if now.Hour() == 0 && now.Minute() == 0 {

			if s.day.every {

				if now.Day()%s.day.value == 0 {

					go s.fn()
				}

			} else {

				//时间区间
				if s.day.between != nil {

					if now.Day() >= s.day.between.min && now.Day() <= s.day.between.max {

						go s.fn()
					}

				} else {

					if now.Day() == s.day.value {

						go s.fn()
					}
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

			//时间区间
			if s.day.between != nil {

				if now.Day() >= s.day.between.min && now.Day() <= s.day.between.max {

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

}
