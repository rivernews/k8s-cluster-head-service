package utilities

import (
	"log"
	"strings"
)

func StringBuilder(stringList ...string) strings.Builder {
	var stringBuilder strings.Builder
	for _, v := range stringList {
		stringBuilder.WriteString(v)
	}
	return stringBuilder
}

func BuildString(stringList ...string) string {
	stringBuilder := StringBuilder(stringList...)
	return stringBuilder.String()
}

func Logger(logLevel string, stringList ...string) {
	value, exist := LogLevelTypes[logLevel]
	if exist && GetLogLevelValue() >= value {
		var prefix string

		if value == LogLevelTypes["DEBUG"] {
			prefix = "🐛 DEBUG: "
		} else if value == LogLevelTypes["INFO"] {
			prefix = "ℹ️ INFO: "
		} else if value == LogLevelTypes["WARN"] {
			prefix = "🟠 WARN: "
		} else if value == LogLevelTypes["ERROR"] {
			prefix = "🛑 ERROR: "
		}

		var logBuilder strings.Builder
		logBuilder.WriteString(prefix)
		for _, v := range stringList {
			logBuilder.WriteString(v)
		}
		log.Println(logBuilder.String())
	}
}
