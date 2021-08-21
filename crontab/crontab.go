package crontab

type crontab struct {
	schedules []*schedule
}

type schedule struct {

	//秒、分、小时、天、月、年，以秒换算
	//

}

func Registered(c *crontab) {

	c.hourly().function(func() {

	})

}

func (c *crontab) hourly() *crontab {

	c.schedules = append(c.schedules, &schedule{})

	return c
}

func (c *crontab) function(fun func()) {

	fun()
}
