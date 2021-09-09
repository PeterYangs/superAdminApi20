package model

type File struct {
	Base
	Id      uint   `json:"id" fillable:"Path,Name,Size,AdminId"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	AdminId int    `json:"admin_id"`
	Admin   Admin  `json:"admin"`
}
