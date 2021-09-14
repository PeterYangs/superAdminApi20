package upload

import (
	"encoding/json"
	"fmt"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/PeterYangs/tools"
	"github.com/PeterYangs/tools/file/read"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"os"
	"time"
)

var upload = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type info struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Type string `json:"type"`
	Nums int    `json:"nums"`
}

func Upload(c *contextPlus.Context) *response.Response {

	form, _ := c.MultipartForm()
	//files := form.File["upload[]"]
	files := form.File["file[]"]

	if len(files) <= 0 {

		return response.Resp().Api(2, "上传文件为空！", "")
	}

	path := make([]string, len(files))

	for i, file := range files {

		//log.Println()

		ex, err := tools.GetExtensionName(file.Filename)

		if err != nil {

			return response.Resp().Api(2, err.Error(), "")
		}

		if !tools.InArray(tools.Explode(",", os.Getenv("ALLOW_UPLOAD_TYPE")), ex) {

			return response.Resp().Api(2, "该拓展类型不允许上传", "")
		}

		date := tools.Date("Ymd", time.Now().Unix())

		os.MkdirAll("uploads/"+date, 0755)

		name := date + "/" + uuid.NewV4().String() + "." + ex

		// 上传文件至指定目录
		c.SaveUploadedFile(file, "uploads/"+name)

		path[i] = name
	}

	if len(path) > 1 {

		return response.Resp().Api(1, "success", path)
	}

	return response.Resp().Api(1, "success", path[0])
}

// BigFile 大文件上传
func BigFile(c *contextPlus.Context) *response.Response {

	conn, err := upload.Upgrade(c.Writer, c.Request, nil)

	if err != nil {

		fmt.Println(err)

		return response.Resp().Api(1, err.Error(), "")
	}

	go func() {

		defer conn.Close()

		tempDir := ""

		var tempListName []string

		var info info

		currentNum := 0

		//清理临时文件
		defer func() {

			for _, s := range tempListName {

				os.Remove("uploads/temp/" + tempDir + "/" + s)

			}

			os.Remove("uploads/temp/" + tempDir)

		}()

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()

			if err != nil {

				fmt.Println(err)

				return
			}

			if msgType == 1 {

				//设置文件信息
				err := json.Unmarshal(msg, &info)

				if err != nil {

					fmt.Println(err)

					return
				}

				tempDir = uuid.NewV4().String()

				//生成临时文件夹
				os.MkdirAll("uploads/temp/"+tempDir, 0755)

			}

			if msgType == 2 {

				//获取文件拓展名
				exName, err := tools.GetExtensionName(info.Name)

				if err != nil {

					fmt.Println(err)

					return
				}

				//过滤非法文件
				if !tools.InArray(tools.Explode(",", os.Getenv("ALLOW_UPLOAD_TYPE")), exName) {

					conn.WriteJSON(map[string]interface{}{"code": 3, "msg": "不允许上传该类型", "data": ""})

					return
				}

				//临时文件名称
				tempName := uuid.NewV4().String() + ".temp"

				f, err := os.OpenFile("uploads/temp/"+tempDir+"/"+tempName, os.O_CREATE|os.O_RDWR, 0644)

				if err != nil {

					fmt.Println(err)

					return
				}

				f.Write(msg)

				f.Close()

				currentNum++

				//发送给客户端进度
				conn.WriteJSON(map[string]interface{}{"code": 2, "msg": "success", "data": currentNum})

				//将临时文件写入到数组
				tempListName = append(tempListName, tempName)

				//当发送到最后一块时
				if currentNum == info.Nums {

					date := tools.Date("Ymd", time.Now().Unix())

					dir := "uploads/" + date

					//生成上传文件夹
					os.MkdirAll(dir, 0775)

					//生成文件名称
					fileName := uuid.NewV4().String() + "." + exName

					//文件完整路径
					path := dir + "/" + fileName

					ff, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0664)

					if err != nil {

						fmt.Println(err)

						return
					}

					//拼接临时文件到目标文件
					for _, s := range tempListName {

						data, err := read.Open("uploads/temp/" + tempDir + "/" + s).Read()

						if err != nil {

							fmt.Println(err)

							ff.Close()

							return
						}

						ff.Write(data)

					}

					//发送文件上传完毕通知
					conn.WriteJSON(map[string]interface{}{"code": 1, "msg": "success", "data": map[string]interface{}{"path": date + "/" + fileName, "name": info.Name, "size": info.Size}})

					//关闭文件流
					ff.Close()

				}

			}

		}

	}()

	return response.Resp().Nil()
}

func checkOrigin(r *http.Request) bool {

	return true
}
