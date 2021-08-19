package upload

import (
	"gin-web/contextPlus"
	"gin-web/response"
	"github.com/PeterYangs/tools"
	uuid "github.com/satori/go.uuid"
	"os"
	"time"
)

func Upload(c *contextPlus.Context) *response.Response {

	form, _ := c.MultipartForm()
	//files := form.File["upload[]"]
	files := form.File["file[]"]

	if len(files) <= 0 {

		return response.Resp().Api(2, "上传文件为空！", "")
	}

	path := make([]string, len(files))

	for i, file := range files {

		//log.Println()

		ex, err := tools.GetExtensionName(file.Filename)

		if err != nil {

			return response.Resp().Api(2, err.Error(), "")
		}

		if !tools.InArray(tools.Explode(",", os.Getenv("ALLOW_UPLOAD_TYPE")), ex) {

			return response.Resp().Api(2, "该拓展类型不允许上传", "")
		}

		date := tools.Date("Ymd", time.Now().Unix())

		os.MkdirAll("uploads/"+date, 0755)

		name := date + "/" + uuid.NewV4().String() + "." + ex

		// 上传文件至指定目录
		c.SaveUploadedFile(file, "uploads/"+name)

		path[i] = name
	}

	if len(path) > 1 {

		return response.Resp().Api(1, "success", path)
	}

	return response.Resp().Api(1, "success", path[0])
}
