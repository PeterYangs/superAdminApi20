package main

import (
	"gin-web/routes"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"runtime"
)

func main() {

	r := gin.Default()

	//加载路由
	routes.Load(r)

	//设置端口
	port := os.Getenv("PORT")

	if port == "" {

		port = "8887"
	}

	//r.Run(":" + port)

	sysType := runtime.GOOS

	//支持平滑重启
	if sysType == "linux" {
		// LINUX系统

		endless.ListenAndServe(":"+port, r)
	}

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
