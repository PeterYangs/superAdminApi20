package file

import (
	"fmt"
	"gin-web/common"
	"gin-web/contextPlus"
	"gin-web/database"
	"gin-web/model"
	"gin-web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"os"
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

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})
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

		return response.Resp().Json(gin.H{"code": 2, "msg": err.Error()})

	}

	var file model.File

	database.GetDb().First(&file, form.Id)

	database.GetDb().Delete(&model.File{}, form.Id)

	//删除对应文件
	err = os.Remove("uploads/" + file.Path)

	fmt.Println(err)
	fmt.Println("uploads/" + file.Path)

	return response.Resp().Api(1, "success", "")

}
