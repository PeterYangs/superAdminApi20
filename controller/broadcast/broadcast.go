package broadcast

import (
	"fmt"
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/response"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var broadcast = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type Online struct {
	Total int64
	List  sync.Map
	Lock  sync.Mutex
}

var Onlines Online

func (o *Online) Add() {

}

func Broadcast(c *contextPlus.Context) *response.Response {

	conn, err := broadcast.Upgrade(c.Writer, c.Request, nil)

	if err != nil {

		fmt.Println(err)

		return response.Resp().Api(1, err.Error(), "")
	}

	go func() {

		defer conn.Close()

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
