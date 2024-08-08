package ws

import (
	"fmt"
	"go.gin.order/src/config/messagequeue"
	"log"
	"sync"
)

type Hub struct {
	clients    map[*wsClient]bool
	register   chan *wsClient
	unregister chan *wsClient
	rooms      map[string]map[*wsClient]bool
	joinRoom   chan *wsClient
	leaveRoom  chan *wsClient
	broadcast  chan Message
	privateMsg chan Message
	mqClient   map[*wsClient]string
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		joinRoom:   make(chan *wsClient),
		leaveRoom:  make(chan *wsClient),
		clients:    make(map[*wsClient]bool),
		rooms:      make(map[string]map[*wsClient]bool),
		mqClient:   make(map[*wsClient]string),
		broadcast:  make(chan Message),
		privateMsg: make(chan Message),
	}
}
func (h *Hub) Run(mq *messagequeue.RabbitMQ) {
	for {
		select {
		case client := <-h.register:
			fmt.Println("注册聊天")
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			fmt.Println("退出聊天")
			if _, ok := h.clients[client]; ok {
				h.mu.Lock()
				delete(h.clients, client)
				close(client.send)
				if _, ok = h.rooms[client.data.room]; ok {
					delete(h.rooms[client.data.room], client)
					if len(h.rooms[client.data.room]) < 1 {
						delete(h.rooms, client.data.room)
					}
				}
				if consumertag, ok := h.mqClient[client]; ok {
					mq.Cancel(consumertag) //中断连接
				}
				h.mu.Unlock()
			}
		case message := <-h.broadcast:
			fmt.Println("公共消息1")
			h.mu.Lock()
			if roomClients, ok := h.rooms[message.Room]; ok {
				fmt.Println("公共消息2")
				for client := range roomClients {
					select {
					case client.send <- []byte(message.Content):
					}
				}
			} else {
				log.Println("公共消息3")
				for client := range h.clients {
					log.Println(client.data.id, "连接实例")
					client.send <- []byte(message.Content)
				}
			}
			h.mu.Unlock()
		case message := <-h.privateMsg:
			h.mu.Lock()
			fmt.Println("私人信息")
			for client := range h.clients {
				if client.data.id == message.Target {
					fmt.Println(message.Target, "cient", client.data.id)
					select {
					case client.send <- []byte(message.Content):
					default:
						close(client.send)
					}
				}
			}
			h.mu.Unlock()
		case client := <-h.joinRoom:
			h.mu.Lock()
			if _, ok := h.rooms[client.data.room]; !ok {
				h.rooms[client.data.room] = make(map[*wsClient]bool)
			}
			h.rooms[client.data.room][client] = true
			h.mu.Unlock()
		case client := <-h.leaveRoom:
			h.mu.Lock()
			if _, ok := h.rooms[client.data.room]; ok {
				delete(h.rooms[client.data.room], client)
				if len(h.rooms[client.data.room]) < 1 {
					delete(h.rooms, client.data.room)
				}
			}
			h.mu.Unlock()

		}
	}
}
