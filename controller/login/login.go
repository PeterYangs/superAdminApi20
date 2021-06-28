package login

import (
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"github.com/gin-gonic/gin"
)

func Login(c *contextPlus.Context) interface{} {

	type Form struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Captcha  string `json:"captcha"  form:"captcha" binding:"required"`
	}

	var form Form

	err := c.ShouldBind(&form)

	if err != nil {

		return gin.H{"code": 2, "mgs": err.Error()}

	}

	if !c.CheckCaptcha(form.Captcha) {

		return gin.H{"code": 2, "mgs": "验证码错误"}
	}

	re := database.GetDb().Where("username = ?", form.Username).Where("password = ?", common.HmacSha256(form.Password)).First(&model.Admin{})

	if re.Error != nil {

		return gin.H{"code": 2, "mgs": "密码错误"}
	}

	return gin.H{"code": 1, "msg": "success"}
}

func Registered(c *contextPlus.Context) interface{} {

	type Validator struct {
		Username   string `json:"username" form:"username" binding:"required"`
		Password   string `json:"password" form:"password" binding:"required"`
		RePassword string `form:"repassword" binding:"required"`
		Email      string `json:"email" form:"email" binding:"required"`
	}

	var form Validator

	err := c.ShouldBind(&form)

	if err != nil {

		return gin.H{"code": 2, "mgs": err.Error()}

	}

	if form.Password != form.RePassword {

		return gin.H{"code": 2, "mgs": "两次密码不一致"}
	}

	ok := database.GetDb().Where("username = ?", form.Username).Or("email = ?", form.Email).First(&model.Admin{})

	if ok.Error == nil {

		return gin.H{"code": 1, "msg": "用户名已被注册", "data": form}
	}

	database.GetDb().Create(&model.Admin{
		Username: form.Username,
		Password: form.Password,
		Email:    form.Email,
	})

	return gin.H{"code": 1, "msg": "success"}

}
