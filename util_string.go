package main

import (
	"reflect"
	"strconv"
	"strings"
)

func trimString(s string) string {
	return strings.Trim(stringlize(s), " 	")
}

func stringlize(obj interface{}) string {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.String:
		return obj.(string)
	case reflect.Int:
		return strconv.Itoa(obj.(int))
	default:
		return ""
	}
}

func string2int(str string) int {
	if len(str) > 0 {
		num, _ := strconv.Atoi(str)
		return num
	}
	return 0
}

func stringToArray(str string, seprate string) []string {
	strSlice := strings.Split(str, ",")
	resultSlice := make([]string, len(strSlice))
	for i, s := range strSlice {
		resultSlice[i] = trimString(s)
	}
	return resultSlice
}

func toCamelCase(s string) string {
	words := strings.Split(s, "_")

	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	result := strings.Join(words, "")

	return result
}
