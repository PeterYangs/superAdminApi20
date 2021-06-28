package model

import (
	"gin-web/types"
	"gorm.io/gorm"
)

type Admin struct {
	Id        uint
	CreatedAt types.Time     `json:"created_at"`
	UpdatedAt types.Time     `json:"updated_at"`
	Username  string         `json:"username" form:"username"`
	Password  string         `json:"password" form:"password" `
	Email     string         `json:"email" form:"email" `
	Status    int            `json:"status" form:"status"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
