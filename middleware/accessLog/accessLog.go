package accessLog

import (
	"bytes"
	"gin-web/contextPlus"
	"gin-web/queue"
	"gin-web/task/access"
	"io/ioutil"
)

func AccessLog(c *contextPlus.Context) {

	//fmt.Println(c.FullPath())
	//

	//cc:=c.Copy()
	////

	b, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	//fmt.Printf("ctx.Request.body: %v", string(data))

	//fmt.Println(string(b))

	//s:=make(map[string]interface{},0)
	//
	//c.ShouldBindBodyWith(&s,binding.JSON)
	//
	//
	//fmt.Println(s)

	//c.PostForm()
	//c.Request.ParseMultipartForm(128)//保存表单缓存的内存大小128M
	//data := c.Request.Form
	//fmt.Println(data)

	//json := make(map[string]interface{},0) //注意该结构接受的内容
	//
	//err:=c.BindJSON(&json)
	//
	//fmt.Println(err)
	//
	//fmt.Println(json)

	//ccc:=c.Copy()
	//
	//b, _ := ccc.GetRawData()
	//
	//fmt.Println(string(b))

	//fmt.Println(c.Request.URL)

	queue.Dispatch(access.NewTask(c.ClientIP(), c.Request.URL.String(), string(b))).Run()

}
