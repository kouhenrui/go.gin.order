package ws

import (
	"fmt"
	"go.gin.order/src/config/messagequeue"
	"log"
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
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		joinRoom:   make(chan *wsClient),
		leaveRoom:  make(chan *wsClient),
		clients:    make(map[*wsClient]bool),
		rooms:      make(map[string]map[*wsClient]bool),
		broadcast:  make(chan Message),
		privateMsg: make(chan Message),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("注册聊天")
			h.clients[client] = true
			go NewConsumerHub(h)
		case client := <-h.unregister:
			fmt.Println("退出聊天")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if _, ok = h.rooms[client.data.room]; ok {
					delete(h.rooms[client.data.room], client)
					if len(h.rooms[client.data.room]) < 1 {
						delete(h.rooms, client.data.room)
					}
				}
			}
		case message := <-h.broadcast:
			fmt.Println("公共消息1")
			if roomClients, ok := h.rooms[message.Room]; ok {
				fmt.Println("公共消息2")
				for client := range roomClients {
					select {
					case client.send <- []byte(message.Content):

						//default:
						//	close(client.send)
						//	delete(h.clients, client)
						//	delete(roomClients, client)
					}
				}
			} else {
				fmt.Println("公共消息3")
				for client := range h.clients {
					client.send <- []byte(message.Content)
				}
			}

		case message := <-h.privateMsg:
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
		case client := <-h.joinRoom:
			if _, ok := h.rooms[client.data.room]; !ok {
				h.rooms[client.data.room] = make(map[*wsClient]bool)
			}
			h.rooms[client.data.room][client] = true

		case client := <-h.leaveRoom:
			if _, ok := h.rooms[client.data.room]; ok {
				delete(h.rooms[client.data.room], client)
				if len(h.rooms[client.data.room]) < 1 {
					delete(h.rooms, client.data.room)
				}
			}

		}
	}
}
func NewConsumerHub(hub *Hub) {
	go func() {
		msgs, err := messagequeue.RabbitChannel.Consume(
			"approval", // queue
			"",         // consumer
			true,       // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
		}

		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			hub.broadcast <- Message{
				Content: string(msg.Body),
				Type:    "msg",
				Target:  "",
				Room:    "",
				Id:      "",
			}
		}
	}()
}
