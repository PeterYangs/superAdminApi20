package main

import (
	"github.com/PeterYangs/superAdminCore/core"
	"superadmin/conf"
	"superadmin/crontab"
	"superadmin/middleware"
	"superadmin/queue"
	"superadmin/routes"
)

func main() {

	c := core.NewCore()

	//加载配置
	c.LoadConf(conf.Conf)

	c.LoadMiddleware(middleware.Load)

	//加载路由
	c.LoadRoute(routes.Routes)

	//加载任务调度
	c.LoadCrontab(crontab.Crontab)

	//加载消息队列
	c.LoadQueues(queue.Queues)

	//启动
	c.Start()

}
