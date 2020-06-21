package utilities

import (
	"strings"
)

func StringBuilder(stringList []string) strings.Builder {
	var stringBuilder strings.Builder
	for _, v := range stringList {
		stringBuilder.WriteString(v)
	}
	return stringBuilder
}

func BuildString(stringList []string) string {
	stringBuilder := StringBuilder(stringList)
	return stringBuilder.String()
}
