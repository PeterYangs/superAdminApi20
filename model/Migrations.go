package model

type Migrations struct {
	Id        uint
	Migration string
	Batch     int
}

//func (Migrations) TableName() string {
//
//	return "migrations"
//}
