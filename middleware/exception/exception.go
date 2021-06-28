package exception

import (
	"fmt"
	"gin-web/contextPlus"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"os"
	"runtime/debug"
	"time"
)

// Exception 错误捕获
func Exception(c *contextPlus.Context) {

	defer func() {
		if r := recover(); r != nil {

			defer c.Abort()

			f, err := os.OpenFile("log/err.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 644)

			if err != nil {

				//c.Abort()
				return
			}

			defer f.Close()

			msg := "[" + tools.Date("Y-m-d H:i:s", time.Now().Unix()) + "]  " + fmt.Sprint(r) + "\n" + string(debug.Stack())

			f.Write([]byte(msg))

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
