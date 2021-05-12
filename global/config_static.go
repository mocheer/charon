package global

// StaticConfig 静态资源配置
type StaticConfig struct {
	// 名称
	Name string
	// 路由模式 history | hash
	Mode string
	// 静态资源根目录
	Dir string
}
