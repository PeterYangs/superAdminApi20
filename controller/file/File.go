package file

import (
	"fmt"
	"gin-web/contextPlus"
	"gin-web/response"
)

func File(c *contextPlus.Context) *response.Response {

	file, _ := c.FormFile("file")

	//file.

	err := c.SaveUploadedFile(file, "./gg.png")

	fmt.Println(err)

	return nil
}
