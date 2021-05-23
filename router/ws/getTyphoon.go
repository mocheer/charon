package ws

// 实时台风信息推送
func getTyphoon() {
	// test
	broadcast <- &msg{Name: "type", Data: "data"}
}
