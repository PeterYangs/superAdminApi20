package crontab

type crontab struct {
	schedules []*schedule
}

type schedule struct {

	//秒、分、小时、天、月、年，以秒换算
	//
	year   int
	month  int
	day    int
	hour   int
	minute int
	second int
}

func Registered(c *crontab) {

	c.hourly().function(func() {

	})

	c.hourlyAt(10).function(func() {

	})

}

func (c *crontab) hourly() *crontab {

	c.schedules = append(c.schedules, &schedule{
		hour: 1,
	})

	return c
}

func (c *crontab) hourlyAt(minute int) *crontab {

	c.schedules = append(c.schedules, &schedule{
		hour:   1,
		minute: minute,
	})

	return c
}

func (c *crontab) function(fun func()) {

	fun()
}
