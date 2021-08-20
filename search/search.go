package search

import (
	"gorm.io/gorm"
)

type Field struct {
	Key       string
	Condition string
}

type search struct {
}

func NewSearch(tx *gorm.DB, maps map[string]interface{}, fields []Field) *search {

	if maps == nil || len(maps) <= 0 {

		return &search{}
	}

	for _, field := range fields {

		if maps[field.Key] != "" {

			switch field.Condition {

			case "like":
				tx.Where(field.Key+" "+field.Condition+"?", "%"+maps[field.Key].(string)+"%")
			default:
				tx.Where(field.Key+" "+field.Condition+"?", maps[field.Key])
			}

		}

	}

	return &search{}

}
