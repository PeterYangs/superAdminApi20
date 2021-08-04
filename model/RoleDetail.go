package model

type RoleDetail struct {
	Base
	Id      uint `json:"id"`
	AdminId int  `json:"admin_id"`
	RoleId  int  `json:"role_id"`
}
