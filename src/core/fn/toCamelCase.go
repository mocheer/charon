package fn

import (
	"regexp"
	"strings"
)

// toCamelCase
func toCamelCase(val string) string {
	r, _ := regexp.Compile("-+(.)?")
	return r.ReplaceAllStringFunc(val, strings.ToUpper)
}
