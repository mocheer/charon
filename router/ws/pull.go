package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

//
func pull() {
	go listen()
}

// 监听交互信息
func listen() {
	// 注册、接收、推送、注销
	for {
		select {
		// 注册
		case connection := <-register:
			clients[connection] = client{}
		// 接收信息
		// 推送信息
		case message := <-broadcast:
			data, err := json.Marshal(message)
			if err == nil {
				for connection := range clients {

					if err := connection.WriteMessage(websocket.TextMessage, data); err != nil {
						log.Println("write error:", err)
						connection.WriteMessage(websocket.CloseMessage, []byte{})
						connection.Close()
						delete(clients, connection)
					}
				}
			}
		// 注销
		case connection := <-unregister:
			delete(clients, connection)
		}
	}
}
