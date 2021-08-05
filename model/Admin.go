package model

import (
	"gin-web/types"
	"gorm.io/gorm"
)

type Admin struct {
	Id        uint           `json:"id"`
	CreatedAt types.Time     `json:"created_at"`
	UpdatedAt types.Time     `json:"updated_at"`
	Username  string         `json:"username" form:"username"`
	Password  string         `json:"-" form:"password" `
	Email     string         `json:"email" form:"email" `
	Status    int            `json:"status" form:"status" gorm:"default:1"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	//RoleId    int            `json:"role_id"`
	RoleDetail RoleDetail `json:"role_detail" gorm:"foreignKey:admin_id"`
}
