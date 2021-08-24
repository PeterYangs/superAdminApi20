package crontab

import (
	"fmt"
	"gin-web/component/logs"
)

type crontab struct {
	schedules []*schedule
}

type schedule struct {

	//秒、分、小时、天、月、年，以秒换算
	//
	//year   int
	month   *number
	day     *number
	hour    *number
	minute  *number
	second  *number
	week    *number //0-6
	crontab *crontab
	fn      func()
	first   bool
}

type number struct {
	every   bool //每
	value   int  //数值
	between *between
}

type between struct {
	min int
	max int
}

func (c *crontab) newSchedule() *schedule {

	return &schedule{
		crontab: c,
		first:   true,
	}
}

//每天
func (s *schedule) everyDay() *schedule {

	if s.first {

		sc := &schedule{
			day: &number{
				every: true,
				value: 1,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.day = &number{
		every: true,
		value: 1,
	}

	return s

}

//某天
func (s *schedule) dayAt(day int) *schedule {

	if s.first {

		sc := &schedule{
			day: &number{
				//every: true,
				value: day,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.day = &number{
		//every: true,
		value: day,
	}

	return s

}

//每几天
func (s *schedule) everyDayAt(day int) *schedule {

	if s.first {

		sc := &schedule{

			day: &number{
				value: day,
				every: true,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.day = &number{
		value: day,
		every: true,
	}

	return s

}

//天，时间区间
func (s *schedule) dayBetween(min, max int) *schedule {

	if s.first {

		sc := &schedule{
			day: &number{
				between: &between{
					min: min,
					max: max,
				},
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.day = &number{
		between: &between{
			min: min,
			max: max,
		},
	}

	return s

}

//每小时
func (s *schedule) everyHour() *schedule {

	if s.first {

		sc := &schedule{
			hour: &number{
				every: true,
				value: 1,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.hour = &number{
		every: true,
		value: 1,
	}

	return s

}

//某一个小时
func (s *schedule) hourlyAt(hour int) *schedule {

	if s.first {

		sc := &schedule{
			hour: &number{
				value: hour,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.hour = &number{
		value: hour,
	}

	return s
}

//每几个小时
func (s *schedule) everyHourAt(hour int) *schedule {

	if s.first {

		sc := &schedule{

			hour: &number{
				value: hour,
				every: true,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.hour = &number{
		value: hour,
		every: true,
	}

	return s

}

//小时，时间区间
func (s *schedule) hourBetween(min, max int) *schedule {

	if s.first {

		sc := &schedule{
			hour: &number{
				between: &between{
					min: min,
					max: max,
				},
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.hour = &number{
		between: &between{
			min: min,
			max: max,
		},
	}

	return s

}

//每分钟
func (s *schedule) everyMinute() *schedule {

	if s.first {

		sc := &schedule{

			minute: &number{
				value: 1,
				every: true,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.minute = &number{
		value: 1,
		every: true,
	}

	return s
}

//每几分钟
func (s *schedule) everyMinuteAt(minute int) *schedule {

	if s.first {

		sc := &schedule{

			minute: &number{
				value: minute,
				every: true,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.minute = &number{
		value: minute,
		every: true,
	}

	return s
}

//某个分钟时间点
func (s *schedule) minuteAt(minute int) *schedule {

	if s.first {

		sc := &schedule{

			minute: &number{
				value: minute,
			},
			first: false,
		}

		s.crontab.schedules = append(s.crontab.schedules, sc)

		return sc

	}

	s.minute = &number{
		value: minute,
	}

	return s

}

func (s *schedule) function(fun func()) {

	//go fun()

	//fmt.Println(fun)

	f := func() {

		//捕获协程异常
		defer func() {

			if r := recover(); r != nil {

				//fmt.Println(r)

				//fmt.Println(debug.Stack())

				msg := fmt.Sprint(r)

				msg = logs.NewLogs().Error(msg).Message()

				fmt.Println(msg)

			}

		}()

		fun()

	}

	s.fn = f
}
