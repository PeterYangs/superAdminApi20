package login

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/database"
	"github.com/PeterYangs/superAdminCore/response"
	regexp "github.com/dlclark/regexp2"
	"strings"
	"superadmin/common"
	"superadmin/model"
)

func Login(c *contextPlus.Context) *response.Response {

	type Form struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
		Captcha  string `json:"captcha"  form:"captcha" binding:"required"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")

	}

	if !c.CheckCaptcha(form.Captcha) {

		return response.Resp().Api(2, "验证码错误", "")
	}

	var admin model.Admin

	//hash,err:=common.Ch

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	re := database.GetDb().Where("username = ?", strings.TrimSpace(form.Username)).First(&admin)

	if re.Error != nil {

		return response.Resp().Api(2, "用户不存在", "")
	}

	hash := admin.Password

	ok := common.CheckPasswordHash(form.Password, hash)

	if !ok {

		return response.Resp().Api(2, "密码错误", "")
	}

	c.Session().Set("admin", admin)

	return response.Resp().Api(1, "success", "")
}

// Registered 后台管理员添加
func Registered(c *contextPlus.Context) *response.Response {

	type Validator struct {
		Username   string `json:"username" form:"username" binding:"required"`
		Password   string `json:"password" form:"password" `
		RePassword string `form:"repassword"`
		Email      string `json:"email" form:"email" binding:"required"`
		RoleId     int    `json:"role_id" form:"role_id"`
		Id         int    `json:"id" form:"id"`
	}

	var form Validator

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")

	}

	//新增
	if form.Id == 0 {

		if form.Password == "" {

			return response.Resp().Api(2, "密码不能为空", "")

		}

		if form.Password != form.RePassword {

			return response.Resp().Api(2, "两次密码不一致", "")
		}

		ok := database.GetDb().Where("username = ?", form.Username).Or("email = ?", form.Email).First(&model.Admin{})

		if ok.Error == nil {

			return response.Resp().Api(2, "用户名已被注册", "")

		}

	} else {

		if form.Password != form.RePassword {

			return response.Resp().Api(2, "两次密码不一致", "")
		}

	}

	if form.Password != "" {

		reg := regexp.MustCompile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,100}$`, 0)

		ok, _ := reg.MatchString(form.Password)

		if !ok {

			return response.Resp().Api(2, "密码必须由数字和大小写字母组成", "")
		}

	}

	tx := database.GetDb().Begin()

	hash, err := common.HashPassword(form.Password)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	admin := model.Admin{
		Username: form.Username,
		Password: hash,
		Email:    form.Email,
		Id:       uint(form.Id),
	}

	var omits []string

	//密码为空则忽略字段更新
	if form.Password == "" {

		omits = append(omits, "password")
	}

	//不能更新用户名
	if form.Id != 0 {

		omits = append(omits, "username")
	}

	err = common.UpdateOrCreateOne(tx, &model.Admin{}, map[string]interface{}{"id": admin.Id}, &admin, omits...)

	if err != nil {

		tx.Rollback()

		return response.Resp().Api(2, err.Error(), "")
	}

	err = common.UpdateOrCreateOne(tx, &model.RoleDetail{}, map[string]interface{}{"admin_id": admin.Id}, &model.RoleDetail{AdminId: int(admin.Id), RoleId: form.RoleId})

	if err != nil {

		tx.Rollback()

		return response.Resp().Api(2, err.Error(), "")
	}

	tx.Commit()

	return response.Resp().Api(1, "success", "")

}

func Logout(c *contextPlus.Context) *response.Response {

	c.Session().Remove("admin")

	return response.Resp().Api(1, "success", "")
}
