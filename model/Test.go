package model

import "gin-web/types"

type Test struct {
	Id        uint
	Name      string
	Array     types.JsonArray
	Map       types.JsonMap
	CreatedAt types.Time
}
