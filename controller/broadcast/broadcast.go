package broadcast

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/gorilla/websocket"
	"net/http"
	"superadmin/util/online"
	"time"
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

	onlineConn := online.NewConn(conn)

	id := on.Add(onlineConn)

	go func() {

		//defer

		defer func() {

			on.Del(id)

			conn.Close()
		}()

		go func() {

			for {

				time.Sleep(1 * time.Second)

				//wErr := conn.WriteJSON(map[string]interface{}{"ping": "ping"})
				wErr := onlineConn.SendMessage("ping")

				if wErr != nil {

					return
				}

			}

		}()

		for {

			_, msg, err := conn.ReadMessage()

			if err != nil {

				fmt.Println(err)

				return
			}

			conn.WriteJSON(map[string]interface{}{"data": msg})

		}

	}()

	return response.Resp().Nil()
}

func checkOrigin(r *http.Request) bool {

	return true
}
