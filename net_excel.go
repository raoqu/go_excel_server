package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

const EXCEL_PATH_PREFIX = "Upload/"
const EXCEL_PATH_POSTFIX = ".xlsx"

type ExcelParseProperties struct {
	File       string
	Sheet      string
	KeyName    string
	ListFields []string
}

type CachedExcelData struct {
	Properties ExcelParseProperties
	ListType   reflect.Type
	Items      []interface{}
}

type CachedExcelDataPool map[string]*CachedExcelData
type ExcelDataTypePool map[string]reflect.Type

var STRING_TYPE = reflect.TypeOf("")

var globalExcelDataCache CachedExcelDataPool = make(CachedExcelDataPool)
var excelDataTypePool ExcelDataTypePool = make(ExcelDataTypePool)

type GetExcelDataRequest struct {
	File    string `json:"file"`
	Sheet   string `json:"sheet"`
	Fields  string `json:"fields"`
	KeyName string `json:"keyName"`
	Key     string `json:"key"`
}

func excelData2TypeList(parser *ExcelParser, data *ExcelData, newStructType reflect.Type) []interface{} {
	list := make([]interface{}, 0, len(*data))
	// traverse through excel data
	for _, row := range *data {

		ptrValue := reflect.New(newStructType)
		instanceValue := reflect.Indirect(ptrValue)

		// update instance fields
		for i := 0; i < newStructType.NumField(); i++ {
			field := newStructType.Field(i)
			value := parser.getValue(row, field.Tag.Get("title"))
			instanceValue.FieldByName(field.Name).SetString(value)
		}
		list = append(list, instanceValue.Interface())
	}
	return list
}

func parseExcelData(properties ExcelParseProperties) ([]interface{}, error) {
	path := EXCEL_PATH_PREFIX + properties.File + EXCEL_PATH_POSTFIX
	parser := NewExcelParser(path)
	excelData, err := parser.loadData(properties.Sheet)

	newType := makeNewType(parser.ColumnTitles)
	items := make([]interface{}, 0)

	if err == nil {
		items = excelData2TypeList(parser, &excelData, newType)
	}

	return items, err
}

func initExcelData(file string, sheet string, fileds string, keyName string) *CachedExcelData {
	cacheKey := file + "_" + stringlize(sheet)
	val, ok := globalExcelDataCache[cacheKey]
	properties := ExcelParseProperties{
		File:       file,
		Sheet:      sheet,
		KeyName:    keyName,
		ListFields: stringToArray(fileds, ","),
	}
	if !ok {
		list, err := parseExcelData(properties)
		if err == nil {
			val = &CachedExcelData{
				Properties: properties,
				ListType:   makeNewType(properties.ListFields),
				Items:      list,
			}
			globalExcelDataCache[cacheKey] = val
		}
	}

	return val
}

func clearCachedExcelData() {
	for k := range globalExcelDataCache {
		delete(globalExcelDataCache, k)
	}

}

func getExcelList(w http.ResponseWriter, r *http.Request) {
	var request GetExcelDataRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		httpResponseError(w, err)
		return
	}

	if len(trimString(request.File)) == 0 || len(trimString(request.Fields)) == 0 {
		httpResponseFail(w, "Illegal params")
		return
	}

	data := initExcelData(request.File, request.Sheet, request.Fields, request.KeyName)
	list := convertToBriefListType(data, data.ListType)

	if data != nil {
		httpResponseObject(w, list)
		return
	}

	httpResponseObject(w, make([]string, 0))
}

func convertToBriefListType(data *CachedExcelData, briefType reflect.Type) []interface{} {
	list := make([]interface{}, 0)
	var str string
	if data != nil {
		for _, item := range data.Items {
			ptrValue := reflect.New(briefType)
			instance := reflect.Indirect(ptrValue)
			val := reflect.ValueOf(item)
			for i := 0; i < briefType.NumField(); i++ {
				field := briefType.Field(i)
				str = val.FieldByName(field.Name).String()
				instance.FieldByName(field.Name).SetString(str)
			}
			list = append(list, instance.Interface())
		}
	}

	return list
}

func getExcelData(w http.ResponseWriter, r *http.Request) {
	var request GetExcelDataRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		httpResponseError(w, err)
		return
	}

	data := initExcelData(request.File, request.Sheet, request.Fields, request.KeyName)

	if data != nil {
		httpResponseObject(w, getItemFromExcelData(data, request.Key))
		return
	}

	httpResponseError(w, err)
}

func getItemFromExcelData(data *CachedExcelData, key string) interface{} {
	if data != nil {
		keyName := data.Properties.KeyName

		for _, item := range data.Items {
			val := reflect.ValueOf(item)
			if val.FieldByName(strings.Title(keyName)).String() == key {
				return val.Interface()
			}
		}
	}

	return nil
}
