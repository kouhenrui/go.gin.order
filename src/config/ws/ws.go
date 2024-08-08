package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.gin.order/src/config/messagequeue"
	"log"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Target  string `json:"target,omitempty"`
	Room    string `json:"room,omitempty"`
	Id      string `json:"id"`
}

// 定义一个 Data 结构体，用于保存用户的信息
type Data struct {
	id       string
	room     string
	msgType  string
	content  []byte
	userList []string
}

// 定义一个 connection 结构体，用于保存每个连接的信息
type wsClient struct {
	hub  *Hub
	ws   *websocket.Conn // WebSocket 连接
	send chan []byte
	data Data
}

func newData() *Data {
	return &Data{
		id:       "",
		room:     "",
		msgType:  "",
		content:  nil,
		userList: nil,
	}
}

// CheckOrigin防止跨站点的请求伪造
var upGrader = &websocket.Upgrader{
	//设置读取写入字节大小
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		//可以添加验证信息
		return true
	},
}

func WsInit(huber *Hub, w http.ResponseWriter, r *http.Request, mq *messagequeue.RabbitMQ) {
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	id := r.URL.Query().Get("id")
	log.Println("id", id)
	data := newData()
	data.id = id
	connect := &wsClient{hub: huber, ws: ws, data: *data, send: make(chan []byte, 1024)}
	connect.hub.register <- connect
	go connect.readPump(mq)
	go connect.writePump()
	//go connect.readrabbitmq()
	//go connect.handleConnect()
}

func (w *wsClient) readPump(mq *messagequeue.RabbitMQ) {
	for {
		_, message, err := w.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Printf("error: %v", err)
			break
		}
		fmt.Println(msg, "msg")
		switch msg.Type {

		case "data":
			w.data.id = msg.Id
		case "chat":
			w.hub.broadcast <- msg
		case "private":
			w.hub.privateMsg <- msg
		case "join":
			w.data.room = msg.Room
			w.hub.joinRoom <- w
			w.hub.broadcast <- Message{Type: "notification", Content: msg.Sender + " has joined the room", Room: msg.Room}
		case "leave":
			w.hub.leaveRoom <- w
			w.hub.broadcast <- Message{Type: "notification", Content: w.data.id + " has left the room", Room: w.data.room}
			w.data.room = ""
		default:
			fmt.Println("Unknown message type:", msg)
		}

		if mq != nil {
			w.hub.mqClient[w] = w.data.id
			queue, err := mq.DeclareQueue(w.data.id)
			if err != nil {
				log.Printf("Failed to declare queue: %v", err)
				return
			}
			err = mq.BindQueue(queue.Name, "", "direct_exchange")
			if err != nil {
				log.Printf("Failed to bind queue: %v", err)
				return
			}

			msgs, err := mq.Consume(queue.Name)
			if err != nil {
				log.Printf("Failed to consume messages: %v", err)
				return
			}
			go func() {
				for msg := range msgs {
					log.Printf("Received a message: %s", msg.Body)
					w.hub.broadcast <- Message{Type: "notification", Content: string(msg.Body)}
					//c.send <- msg.Body
				}
			}()

		}

	}

}

func (w *wsClient) writePump() {
	//ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message, ok := <-w.send:
			//w.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				w.ws.WriteMessage(websocket.CloseMessage, []byte{})
			}
			w.ws.WriteMessage(websocket.TextMessage, message)
			//wss, err := w.ws.NextWriter(websocket.TextMessage)
			//if err != nil {
			//	return
			//}
			//wss.Write(message)
			//case <-ticker.C:
			//	w.ws.SetWriteDeadline(time.Now().Add(writeWait))
			//	if err := w.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
			//		return
			//	}
		}
	}
}

//func (w *wsClient) NewConsumerHub() {
//	queueName := w.data.id
//	log.Println(queueName, "queuename")
//	mq := messagequeue.NewRabbitMQ()
//	queue, err := mq.DeclareQueue(queueName)
//	if err != nil {
//		log.Fatalf("Could not declare queue: %v", err)
//	}
//	var exchangeName string
//	exchangeName = "direct_exchange"
//	//switch w.data.msgType {
//	//case "direct":
//	//	exchangeName = "direct_exchange"
//	//case "broadcast":
//	//	exchangeName = "fanout_exchange"
//	//case "topic":
//	//	exchangeName = "topic_exchange"
//	//default:
//	//	break
//	//}
//	err = mq.BindQueue(queue.Name, "", exchangeName)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	msgs, err := mq.Consume(queue.Name)
//	if err != nil {
//		log.Fatalf("Failed to register a consumer: %v", err)
//	}
//
//	go func() {
//		for msg := range msgs {
//			log.Printf("Received a message: %s", msg.Body)
//			w.hub.broadcast <- Message{Type: "notification", Content: string(msg.Body)}
//			//c.send <- msg.Body
//		}
//	}()
//}
