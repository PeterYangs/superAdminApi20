package controller

import (
	"gin-web/database"
	"gin-web/model"
	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) gin.H {

	var test model.Test

	database.GetDb().First(&test)

	return gin.H{"code": 1, "msg": "success", "data": test}

}
