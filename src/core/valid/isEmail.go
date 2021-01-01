package valid

import (
	"regexp"
)

// IsEmail 判断字符串是否符合邮箱地址
func IsEmail(value string) bool {
	pattern := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return pattern.MatchString(value)
}
