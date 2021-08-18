package model

type Access struct {
	Base
	Id      uint    `json:"id" fillable:"Ip,Url,Params,AdminId"`
	Ip      string  `json:"ip"`
	Url     string  `json:"url"`
	Params  string  `json:"params"`
	AdminId float64 `json:"admin_id"`
}
