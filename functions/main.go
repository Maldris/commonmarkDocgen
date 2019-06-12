package functions

import (
	"text/template"

	"github.com/Maldris/commonmarkDocgen/functions/num2words"
)

func GetFunctionMap(exts template.FuncMap) template.FuncMap {
	funcMap := template.FuncMap{
		"eq":       eq,
		"cell":     newCell,
		"lower":    lower,
		"upper":    upper,
		"caps":     upper,
		"title":    title,
		"currency": currencyFormat,
		"date":     formatDate,
		"intDate":  integerDateFormat,
		"num2word": num2words.Convert,
		"add":      add,
		"subtract": subtract,
		"multiply": multiply,
		"divide":   divide,
		"cb":       codeBlock,
		"dict":     dictionary,
		"yn":       boolString,
	}

	for key, val := range exts {
		funcMap[key] = val
	}
	return funcMap
}
