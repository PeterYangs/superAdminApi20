package regex

import (
	"gin-web/contextPlus"
	"github.com/gin-gonic/gin"
)

// Regex 参数规则验证示例，路由为 /regex/:name ,请求为 /regex/1sds?test[]=1&test[]=2,regex标记只支持string和[]string两个类型
func Regex(c *contextPlus.Context) interface{} {

	type regex struct {
		Test []string `form:"test[]" json:"test" regex:"[0-9a-z/]+"`
		Name string   `uri:"name" json:"name" regex:"[0-9a-z]+"`
	}

	var t regex

	err := c.ShouldBindPlus(&t)

	if err != nil {

		return gin.H{"code": 2, "mgs": err.Error()}

	}

	return gin.H{"code": 1, "msg": "hello world"}
}
