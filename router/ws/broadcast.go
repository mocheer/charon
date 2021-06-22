package ws

type msg struct {
	Name string `json:"type"`
	Data string `json:"data"`
}

// 广播信息
var broadcast = make(chan *msg)
