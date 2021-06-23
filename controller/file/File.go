package file

import (
	"fmt"
	"gin-web/contextPlus"
)

func File(c *contextPlus.Context) interface{} {

	file, _ := c.FormFile("file")

	//file.

	err := c.SaveUploadedFile(file, "./gg.png")

	fmt.Println(err)

	return nil
}
