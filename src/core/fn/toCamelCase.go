package fn

import (
	"regexp"
	"strings"
)

// ToCamelCase
// fn.ToCamelCase("camel-case") == "camelCase"
func ToCamelCase(val string) string {
	r, _ := regexp.Compile("-+(.)?")
	return r.ReplaceAllStringFunc(val, strings.ToUpper)
}
