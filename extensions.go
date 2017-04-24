package docgen

import (
	"errors"
	"strconv"
	"time"
)

var dateFormat = "2 January 2006"

func currencyFormat(amount float64) string {
	str := "$" + strconv.FormatFloat(amount, 'f', 2, 64)

	return str
}

func formatDate(tim time.Time) string {
	return tim.Format(dateFormat)
}

func integerDateFormat(tim int64) string {
	date := time.Unix(tim, 0).Local()
	return formatDate(date)
}

func add(a, b int) int {
	return a + b
}

func subtract(a, b int) int {
	return a - b
}

func multiply(a, b float64) float64 {
	return a * b
}

func divide(a, b float64) float64 {
	return a / b
}

func codeBlock() string {
	return "```"
}

func dictionary(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func boolString(b bool) string {
	if b {
		return "Y"
	}
	return "N"
}
