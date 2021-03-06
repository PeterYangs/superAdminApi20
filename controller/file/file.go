package file

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/database"
	"github.com/PeterYangs/superAdminCore/v2/response"
	"github.com/spf13/cast"
	"os"
	"superadmin/common"
	"superadmin/model"
)

func Update(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id   int    `json:"id" form:"id"`
		Path string `json:"path" form:"path" binding:"required"`
		Name string `json:"name" form:"name" binding:"required"`
		Size int64  `json:"size" form:"size" binding:"required"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	file := model.File{
		Id:      uint(form.Id),
		Path:    form.Path,
		Name:    form.Name,
		Size:    form.Size,
		AdminId: c.GetAdminId(),
	}

	err = common.UpdateOrCreateOne(database.GetDb(), &model.File{}, map[string]interface{}{"id": form.Id}, &file)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")
	}

	return response.Resp().Api(1, "success", "")
}

func List(c *contextPlus.Context) *response.Response {

	files := make([]*model.File, 0)

	tx := database.GetDb().Preload("Admin").Model(&model.File{}).Order("id desc")

	data := common.Paginate(tx, &files, cast.ToInt(c.DefaultQuery("p", "1")), 10)

	return response.Resp().Api(1, "success", data)
}

func Destroy(c *contextPlus.Context) *response.Response {

	type Form struct {
		Id int `json:"id" uri:"id"`
	}

	var form Form

	err := c.ShouldBindPlus(&form)

	if err != nil {

		return response.Resp().Api(2, err.Error(), "")

	}

	var file model.File

	database.GetDb().First(&file, form.Id)

	database.GetDb().Delete(&model.File{}, form.Id)

	//删除对应文件
	os.Remove("uploads/" + file.Path)

	return response.Resp().Api(1, "success", "")

}
