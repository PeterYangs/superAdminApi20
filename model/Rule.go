package model

type Rule struct {
	Base
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Rule      string `json:"rule"`
	GroupName string `json:"group_name"`
}
