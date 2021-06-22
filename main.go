package main

import (
	"gin-web/common"
	"gin-web/conf"
	"gin-web/middleware/session"
	"gin-web/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"runtime"
)

func main() {

	r := gin.Default()

	//加载配置
	conf.Load()

	//开启session(全局中间件)
	common.GlobalMiddleware(r, session.StartSession)

	//加载路由
	routes.Load(r)

	//设置端口
	port := os.Getenv("PORT")

	if port == "" {

		port = "8887"
	}

	sysType := runtime.GOOS

	////支持平滑重启，kill -1 pid
	//if sysType == "linux" {
	//	// LINUX系统
	//
	//	endless.ListenAndServe(":"+port, r)
	//}

	//windows只做开发测试
	if sysType == "windows" {
		// windows系统

		r.Run(":" + port)

	}

}

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

}
