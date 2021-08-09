package model

type Category struct {
	Base
	Id    uint   `json:"id" fillable:"Pid,Lv,Title,Img,Sort,Path"`
	Pid   int    `json:"pid"`
	Lv    int    `json:"lv"`
	Title string `json:"title"`
	Img   string `json:"img"`
	Sort  int    `json:"sort"`
	Path  string `json:"path"`
}
