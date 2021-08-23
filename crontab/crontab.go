package crontab

import "fmt"

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
}

type number struct {
	every bool //每
	value int  //数值
}

func Registered(c *crontab) {

	c.newSchedule().hourly().function(func() {

	})

	c.newSchedule().hourlyAt(10).function(func() {

	})

	c.newSchedule().everyMinute().function(func() {

		fmt.Println("每分钟")

	})

	c.newSchedule().everyMinuteAt(2).function(func() {

		fmt.Println("每2分钟")

	})

}

func (c *crontab) newSchedule() *schedule {

	return &schedule{
		crontab: c,
	}
}

func (s *schedule) hourly() *schedule {

	sc := &schedule{
		hour: &number{
			every: true,
		},
	}

	s.crontab.schedules = append(s.crontab.schedules, sc)

	return sc
}

func (s *schedule) hourlyAt(minute int) *schedule {

	sc := &schedule{
		hour: &number{
			every: true,
		},
		minute: &number{
			value: minute,
		},
	}

	s.crontab.schedules = append(s.crontab.schedules, sc)

	return sc
}

//每分钟
func (s *schedule) everyMinute() *schedule {

	sc := &schedule{

		minute: &number{
			value: 1,
			every: true,
		},
	}

	s.crontab.schedules = append(s.crontab.schedules, sc)

	return sc

}

//每几分钟
func (s *schedule) everyMinuteAt(minute int) *schedule {

	sc := &schedule{

		minute: &number{
			value: minute,
			every: true,
		},
	}

	s.crontab.schedules = append(s.crontab.schedules, sc)

	return sc

}

func (s *schedule) function(fun func()) {

	//go fun()

	//fmt.Println(fun)

	s.fn = fun
}
