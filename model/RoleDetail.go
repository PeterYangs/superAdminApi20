package model

type RoleDetail struct {
	Base
	Id      uint `json:"id" fillable:"AdminId,RoleId"`
	AdminId int  `json:"admin_id"`
	RoleId  int  `json:"role_id"`
	Role    Role `json:"role"`
}
