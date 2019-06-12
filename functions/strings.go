package functions

import (
	"strconv"
	"strings"
)

func lower(str string) string {
	return strings.ToLower(str)
}

func upper(str string) string {
	return strings.ToUpper(str)
}

func title(str string) string {
	return strings.ToTitle(strings.ToLower(str))
}

func boolString(b bool) string {
	return strconv.FormatBool(b)
}
