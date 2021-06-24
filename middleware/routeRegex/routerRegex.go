package routeRegex

import (
	"gin-web/contextPlus"
)

func RouterRegex(c *contextPlus.Context) {

	//fmt.Println(c.Regex)

	//fmt.Println(c.Handler)

	regex := c.Handler.Regex

	_ = regex

}
