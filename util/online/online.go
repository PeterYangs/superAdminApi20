package online

import (
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type online struct {
	total int64
	list  sync.Map
	lock  sync.Mutex
}

var Online *online

var once sync.Once

type Conn struct {
	id        string
	conn      *websocket.Conn
	lastReply time.Time //上一次回复时间
}

func NewOnline() *online {

	once.Do(func() {

		Online = &online{
			total: 0,
			list:  sync.Map{},
			lock:  sync.Mutex{},
		}

	})

	go Online.checkTime()

	return Online
}

func (o *online) Add(c *Conn) string {

	o.total++

	o.list.Store(c.id, c)

	return c.id

}

func (o *online) Del(id string) {

	o.total--

	o.list.Delete(id)

}

// Total 当前在线人数
func (o *online) Total() int64 {

	return o.total
}

func (o *online) GetConnById(id string) (bool, *Conn) {

	c, ok := o.list.Load(id)

	if ok {

		con := c.(*Conn)

		return ok, con

	}

	return ok, nil
}

func (o *online) checkTime() {

	for {

		time.Sleep(30 * time.Second)

		o.list.Range(func(key, value interface{}) bool {

			con := value.(*Conn)

			if time.Now().Sub(con.lastReply).Seconds() > 60 {

				fmt.Println("超时")

				con.conn.Close()
			}

			return true
		})

	}

}

//------------------------------------------------------------------------------------

func NewConn(conn *websocket.Conn) *Conn {

	return &Conn{id: uuid.NewV4().String(), conn: conn, lastReply: time.Now()}
}

// SendMessage 发送字符串消息
func (c *Conn) SendMessage(message string) error {

	return c.conn.WriteJSON(map[string]interface{}{"type": "message", "content": message})
}
