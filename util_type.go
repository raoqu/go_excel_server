package main

import (
	"reflect"
	"strings"
)

// Make a new type from a string array
//
// e.g. {"a", "b", "c", "d} -->
//
//	struct {
//		  A string `json:"a"`
//		  B string `json:"b"`
//		  C string `json:"c"`
//		  D string `json:"d"`
//		}
func makeNewType(titles []string) reflect.Type {
	fields := []reflect.StructField{}
	for _, val := range titles {
		title := trimString(val)
		camelTitle := toCamelCase(title)
		if len(title) > 0 {
			fields = append(fields, reflect.StructField{
				Name: strings.Title(camelTitle),
				Type: STRING_TYPE,
				Tag:  reflect.StructTag(`json:"` + camelTitle + `" title:"` + title + `"`),
			})
		}
	}

	newType := reflect.StructOf(fields) // make a struct type
	return newType
}
