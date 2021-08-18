package model

type Access struct {
	Base
	Id     uint   `json:"id" fillable:"Ip,Url,Params"`
	Ip     string `json:"ip"`
	Url    string `json:"url"`
	Params string `json:"params"`
}
