package crontab

import (
	"fmt"
	"runtime/debug"
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

func Registered(c *crontab) {

	c.newSchedule().everyHour().function(func() {

		fmt.Println("每小时")

	})

	c.newSchedule().hourlyAt(16).everyMinute().function(func() {

		fmt.Println("每个16点的每分钟")

	})

	c.newSchedule().minuteAt(18).function(func() {

		fmt.Println("每小时的第18分钟")

	})

	c.newSchedule().everyMinute().function(func() {

		panic("模拟报错")

		fmt.Println("每分钟")

	})

	c.newSchedule().everyMinuteAt(2).function(func() {

		fmt.Println("每2分钟")

	})

	c.newSchedule().everyDay().hourlyAt(16).minuteAt(36).function(func() {

		fmt.Println("每天16点36分")

	})

	c.newSchedule().dayAt(23).hourlyAt(16).minuteAt(50).function(func() {

		fmt.Println("23号16点50分")

	})

	c.newSchedule().dayAt(24).hourBetween(8, 10).function(func() {

		fmt.Println("24号8点-10点")

	})

	c.newSchedule().hourBetween(8, 9).everyMinute().function(func() {

		fmt.Println("24号8点-9点每分钟")

	})

	c.newSchedule().dayBetween(22, 24).everyHour().everyMinute().function(func() {

		fmt.Println("22号-24号每分钟")

	})

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

				fmt.Println(r)

				fmt.Println(debug.Stack())

			}

		}()

		fun()

	}

	s.fn = f
}
