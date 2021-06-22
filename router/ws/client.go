package ws

import "github.com/gofiber/websocket/v2"

type client struct {
	// authed bool
	// events []string
}

// 连接集
var clients = make(map[*websocket.Conn]client)

// 注册
var register = make(chan *websocket.Conn)

// 注销
var unregister = make(chan *websocket.Conn)
