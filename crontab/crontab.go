package crontab

import (
	"fmt"
	"gin-web/component/logs"
	"sync"
)

type crontab struct {
	schedules []*schedule
	quitWait  *sync.WaitGroup
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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

	}

	s.day = &number{

		value: day,
	}

	return s

}

//每几天
func (s *schedule) everyDayAt(day int) *schedule {

	if s.first {

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

	}

	s.hour = &number{
		value: hour,
	}

	return s
}

//每几个小时
func (s *schedule) everyHourAt(hour int) *schedule {

	if s.first {

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

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

		s.first = false

		s.crontab.schedules = append(s.crontab.schedules, s)

	}

	s.minute = &number{
		value: minute,
	}

	return s

}

func (s *schedule) function(fun func()) {

	f := func() {

		//定时任务安全退出
		s.crontab.quitWait.Add(1)

		//捕获协程异常
		defer func() {

			if r := recover(); r != nil {

				msg := fmt.Sprint(r)

				msg = logs.NewLogs().Error(msg).Message()

				fmt.Println(msg)

			}

			s.crontab.quitWait.Done()

		}()

		fun()

	}

	s.fn = f
}
