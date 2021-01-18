package util

import (
	"strings"
)

func GetStringInBetweenTwoString(str string, startS string, endS string) string {
	var result string = ""

	s := strings.Index(str, startS)
	if s == -1 {
		return result
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result
	}
	result = newS[:e]
	return result
}

func CopyMap(m map[string]string) map[string]string {
	cp := make(map[string]string)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}
