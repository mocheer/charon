package auth

// LoginInput 请求输入参数
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
