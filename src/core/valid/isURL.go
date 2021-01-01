package valid

import (
	"strings"
)

// 判断字符串是否是一个url地址
func IsURL(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}
