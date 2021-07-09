package exception

import (
	"fmt"
	"gin-web/component/logs"
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
	"os"
)

// Exception 错误捕获
func Exception(c *contextPlus.Context) {

	defer func() {
		if r := recover(); r != nil {

			defer c.Abort()

			msg := fmt.Sprint(r)

			msg = logs.NewLogs().Error(msg).Message()

			if os.Getenv("APP_DEBUG") == "true" {

				c.String(500, msg)

			} else {

				c.JSON(500, gin.H{"code": 500})
			}

			//c.Abort()

		}

	}()

	c.Next()

}
