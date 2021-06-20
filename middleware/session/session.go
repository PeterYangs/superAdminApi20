package session

import (
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"os"
	"strconv"
)

func StartSession(c *gin.Context) {

	//uuid:=uuid

	u := uuid.NewV4().String() + "-" + tools.Md5(os.Getenv("APP_NAME"))

	life, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))

	c.SetCookie("gin_web", u, life, "/", c.Request.Host, false, true)

	c.Next()

}
