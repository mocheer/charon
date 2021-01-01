package valid

import (
	"regexp"
)

// IsNumeric 判断字符串是否是数字
func IsNumeric(str string) bool {
	pattern := regexp.MustCompile(`^[-+]?\d*\.?\d+(?:[eE][-+]?\d+)?$`)
	return pattern.MatchString(str)
}
