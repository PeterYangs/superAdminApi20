package online

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//所有连接处理结构体
type online struct {
	total int64    //在线连接数
	list  sync.Map //所有连接存放
	lock  sync.Mutex
}

var Online *online

var once sync.Once

// Conn 单个连接结构体
type Conn struct {
	id        string          //连接id
	conn      *websocket.Conn //websocket连接对象
	lastReply time.Time       //上一次回复时间
	adminId   int             //账号id
}

type Message struct {
	Code    int         `json:"code"`
	Types   string      `json:"types"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewOnline() *online {

	once.Do(func() {

		Online = &online{
			total: 0,
			list:  sync.Map{},
			lock:  sync.Mutex{},
		}

		go Online.checkTime()

		go Online.pushTotal()
	})

	return Online
}

// Add 添加一个连接
func (o *online) Add(c *Conn) string {

	o.total++

	o.list.Store(c.id, c)

	return c.id

}

// Del 删除一个连接
func (o *online) Del(id string) {

	o.total--

	o.list.Delete(id)

}

// Total 当前在线人数
func (o *online) Total() int64 {

	return o.total
}

// GetConnById 根据id获取连接
func (o *online) GetConnById(id string) (bool, *Conn) {

	c, ok := o.list.Load(id)

	if ok {

		con := c.(*Conn)

		return ok, con

	}

	return ok, nil
}

// SendAllMessage 群发
func (o *online) SendAllMessage(message Message) {

	o.list.Range(func(key, value interface{}) bool {

		con := value.(*Conn)

		con.SendMessage(message)

		return true
	})
}

//心跳检测
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

//群发在线连接数
func (o *online) pushTotal() {

	for {

		time.Sleep(15 * time.Second)

		o.list.Range(func(key, value interface{}) bool {

			con := value.(*Conn)

			con.SendJson(1, "total", "success", o.total)

			return true
		})

	}

}

// GetTotal 获取在线人数
func (o *online) GetTotal() int64 {

	return o.total
}

//------------------------------------------------------------------------------------

func NewConn(conn *websocket.Conn) *Conn {

	return &Conn{id: uuid.NewV4().String(), conn: conn, lastReply: time.Now()}
}

// SetAdminId 设置账号id
func (c *Conn) SetAdminId(adminId int) {

	c.adminId = adminId

}

// SendMessage 结构体式
func (c *Conn) SendMessage(message Message) error {

	return c.conn.WriteJSON(message)
}

// SendJson 函数式
func (c *Conn) SendJson(code int, types string, message string, data interface{}) error {

	return c.conn.WriteJSON(Message{Code: code, Types: types, Message: message, Data: data})
	//return c.conn.WriteJSON(map[string]interface{}{"code": code, "type": types, "message": message, "data": data})
}

// SetReplyTime 设置上一次回复时间
func (c *Conn) SetReplyTime() {

	c.lastReply = time.Now()
}

//---------------------------------------------------------------------------------------

func NewMessage(message []byte) Message {

	var m Message

	_ = json.Unmarshal(message, &m)

	return m

}
