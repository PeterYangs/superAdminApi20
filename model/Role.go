package model

import "superadmin/types"

type Role struct {
	Base
	Id    uint             `json:"id"`
	Title string           `json:"title"`
	Rules types.CommaArray `json:"rules"`
}
