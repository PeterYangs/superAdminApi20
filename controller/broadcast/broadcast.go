package broadcast

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/gorilla/websocket"
	"net/http"
	"superadmin/util/online"
)

var broadcast = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func Broadcast(c *contextPlus.Context) *response.Response {

	conn, err := broadcast.Upgrade(c.Writer, c.Request, nil)

	if err != nil {

		fmt.Println(err)

		return response.Resp().Api(1, err.Error(), "")
	}

	on := online.NewOnline()

	//实例化一个新连接对象
	onlineConn := online.NewConn(conn)

	onlineConn.SetAdminId(c.GetAdminId())

	//添加一个新连接
	id := on.Add(onlineConn)

	//上线后发送一个在线人数数据
	onlineConn.SendJson(1, "total", "success", on.GetTotal())

	go func() {

		//defer

		defer func() {

			on.Del(id)

			conn.Close()
		}()

		for {

			_, msg, err := conn.ReadMessage()

			//msg.

			if err != nil {

				fmt.Println(err)

				return
			}

			//重置上次访问时间
			onlineConn.SetReplyTime()

			message := online.NewMessage(msg)

			switch message.Types {

			case "ping":

				onlineConn.SendJson(1, "ping", "", []string{})

			}

			//conn.WriteJSON(map[string]interface{}{"data": msg})

		}

	}()

	return response.Resp().Nil()
}

// GroupMessage 群发测试
func GroupMessage(c *contextPlus.Context) *response.Response {

	on := online.NewOnline()

	on.SendAllMessage(online.Message{Code: 1, Types: "group_test", Message: "群发测试"})

	return response.Resp().Api(1, "success", []string{})
}

func checkOrigin(r *http.Request) bool {

	return true
}
