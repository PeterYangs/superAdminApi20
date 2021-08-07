package model

type Menu struct {
	Base
	Id    uint   `json:"id" fillable:"Pid,Title,Path,Sort"`
	Pid   int    `json:"pid"`
	Title string `json:"title"`
	Path  string `json:"path"`
	Sort  int    `json:"sort"`
}
