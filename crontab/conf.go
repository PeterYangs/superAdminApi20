package crontab

import (
	"fmt"
)

func Registered(c *crontab) {

	//fmt.Println(c.quitWait)

	//c.newSchedule().everyHour().function(func() {
	//
	//	fmt.Println("每小时")
	//
	//})
	//
	//c.newSchedule().hourlyAt(16).everyMinute().function(func() {
	//
	//	fmt.Println("每个16点的每分钟")
	//
	//})
	//
	//c.newSchedule().minuteAt(18).function(func() {
	//
	//	fmt.Println("每小时的第18分钟")
	//
	//})
	//
	c.newSchedule().everyMinute().function(func() {

		//panic("模拟报错")
		fmt.Println("每分钟")
		//time.Sleep(5 * time.Second)
		//fmt.Println("结束")

	})
	//
	c.newSchedule().everyMinuteAt(2).function(func() {

		fmt.Println("每2分钟")

	})
	//
	//c.newSchedule().everyDay().hourlyAt(16).minuteAt(36).function(func() {
	//
	//	fmt.Println("每天16点36分")
	//
	//})
	//
	//c.newSchedule().dayAt(23).hourlyAt(16).minuteAt(50).function(func() {
	//
	//	fmt.Println("23号16点50分")
	//
	//})
	//
	//c.newSchedule().dayAt(24).hourBetween(8, 10).function(func() {
	//
	//	fmt.Println("24号8点-10点")
	//
	//})
	//
	//c.newSchedule().hourBetween(8, 9).everyMinute().function(func() {
	//
	//	fmt.Println("24号8点-9点每分钟")
	//
	//})
	//
	//c.newSchedule().dayBetween(22, 24).everyHour().everyMinute().function(func() {
	//
	//	fmt.Println("22号-24号每分钟")
	//
	//})

}
