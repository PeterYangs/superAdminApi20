package model

import "gin-web/types"

type Base struct {
	CreatedAt types.Time `json:"created_at"`
	UpdatedAt types.Time `json:"updated_at"`
}
