package bettermolang

import (
	"strings"
)

func removeWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}

func removeComments(input string) string {
	var returnString string = ""
	commentStart := -1
	isString := false
	for i := 0; i < len(input); i++ {
		if input[i] == '"' {
			isString = !isString
		}
		if !isString && input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			commentStart = i
		}
		if commentStart != -1 && input[i] == '\n' {
			commentStart = -1
		}
		if commentStart == -1 {
			returnString += string(input[i])
		}
	}
	return returnString
}
