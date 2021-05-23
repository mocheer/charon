package ws

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Use 初始化 uploadFile 路由
func Use(api fiber.Router) {
	router := api.Group("/ws")
	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	router.Get("/", websocket.New(handle))
	//
	pull()
	push()
}

func handle(c *websocket.Conn) {
	// 注销
	defer func() {
		unregister <- c
		c.Close()
	}()

	// 注册
	register <- c
	//
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			fmt.Println(string(message))
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}

}
