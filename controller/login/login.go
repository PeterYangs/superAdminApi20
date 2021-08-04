package login

import (
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func Login(c *contextPlus.Context) *response.Response {

	type Form struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Captcha  string `json:"captcha"  form:"captcha" binding:"required"`
	}

	var form Form

	err := c.ShouldBind(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "mgs": err.Error()})

	}

	if !c.CheckCaptcha(form.Captcha) {

		return response.Resp().Json(gin.H{"code": 2, "mgs": "验证码错误"})
	}

	var admin model.Admin

	re := database.GetDb().Where("username = ?", form.Username).Where("password = ?", common.HmacSha256(form.Password)).First(&admin)

	if re.Error != nil {

		return response.Resp().Json(gin.H{"code": 2, "mgs": "密码错误"})
	}

	c.Session().Set("admin", admin)

	return response.Resp().Json(gin.H{"code": 1, "mgs": "success"})
}

func Registered(c *contextPlus.Context) *response.Response {

	type Validator struct {
		Username   string `json:"username" form:"username" binding:"required"`
		Password   string `json:"password" form:"password" `
		RePassword string `form:"repassword"`
		Email      string `json:"email" form:"email" binding:"required"`
		RoleId     int    `json:"role_id" form:"role_id" binding:"required"`
		Id         int    `json:"id" form:"id"`
	}

	var form Validator

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	if form.Id == 0 {

		if form.Password != form.RePassword {

			return response.Resp().Json(gin.H{"code": 2, "msg": "两次密码不一致"})
		}

		ok := database.GetDb().Where("username = ?", form.Username).Or("email = ?", form.Email).First(&model.Admin{})

		if ok.Error == nil {

			return response.Resp().Json(gin.H{"code": 2, "msg": "用户名已被注册"})

		}

	} else {

		if form.Password != "" && form.Password != form.RePassword {

			return response.Resp().Json(gin.H{"code": 2, "msg": "两次密码不一致"})
		}

	}

	tx := database.GetDb().Begin()

	admin := model.Admin{
		Username: form.Username,
		Password: common.HmacSha256(form.Password),
		Email:    form.Email,
	}

	updateColumns := []string{"username", "email"}

	if form.Password != "" {

		updateColumns = append(updateColumns, "password")
	}

	re := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&admin)

	if re.Error != nil {

		tx.Rollback()

		return response.Resp().Api(2, re.Error.Error(), "")
	}

	re = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "admin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role_id"}),
	}).Create(&model.RoleDetail{AdminId: int(admin.Id), RoleId: form.RoleId})

	if re.Error != nil {

		tx.Rollback()

		return response.Resp().Api(2, re.Error.Error(), "")
	}

	tx.Commit()

	return response.Resp().Json(gin.H{"code": 1, "mgs": "success"})

}

func Logout(c *contextPlus.Context) *response.Response {

	c.Session().Remove("admin")

	return response.Resp().Json(gin.H{"code": 1, "mgs": "success"})
}
