package upload

import (
	"gin-web/contextPlus"
	"gin-web/response"
	uuid "github.com/satori/go.uuid"
	"log"
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
		log.Println(file.Filename)

		name := "uploads/" + uuid.NewV4().String() + ".png"

		// 上传文件至指定目录
		c.SaveUploadedFile(file, name)

		path[i] = name
	}

	if len(path) > 1 {

		return response.Resp().Api(1, "success", path)
	}

	return response.Resp().Api(1, "success", path[0])
}
