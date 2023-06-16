package main

import (
	"reflect"
	"strconv"

	"github.com/tealeg/xlsx"
)

type RowData []interface{}              // 行数据(多列)
type ExcelData []RowData                // excel 多行数据（不包含标题）
type AliasGroup []string                // 相同含义的列别名
type AliasConfig []AliasGroup           // 所有列别名组
type AliasIndexes map[string]AliasGroup // 列名->别名组
type CustomParseTitleFunc func(values []interface{}) []string

type ExcelParseOption struct {
	ColumnsAliases   AliasConfig          // 标题别名配置 [ ['A', 'field1', '字段一'], ['B', 'field2', '字段二'], ... ]
	CustomParseTitle CustomParseTitleFunc // 自定义title解析器
	TitleRowCount    int                  // 非数据行数
	TitleRowIndex    int                  // 标题所在行
	InvalidColCount  int                  // 无效行数量，默认为1
}

type ExcelParser struct {
	ColumnTitles []string
	Data         ExcelData      // 数据
	TitleData    ExcelData      //
	TitleAliases AliasConfig    // [ ['A', 'field1', '字段一'], ['B', 'field2', '字段二'], ... ]
	TitleIndexes map[string]int // { 'A': 1, 'field1': 1, 'field2': 2 }

	Options ExcelParseOption

	// init params
	Filepath string
}

var charArr = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// 默认解析配置
var DefaultExcelParserOptions = ExcelParseOption{
	ColumnsAliases:   AliasConfig{},
	CustomParseTitle: nil,
	TitleRowCount:    1,
	TitleRowIndex:    1,
	InvalidColCount:  0,
}

func NewExcelParser(path string) *ExcelParser {
	excel := &ExcelParser{
		Data:     make(ExcelData, 0),
		Filepath: path,
	}

	// 初始化默认配置
	excel.Options = DefaultExcelParserOptions

	return excel
}

func (excel *ExcelParser) columnIndexToName(n int) string {
	i := n - 1
	a := ""
	b := charArr[i%26]

	if n > 26 {
		a = charArr[(i-26)/26]
	}
	return a + b
}

func (excel *ExcelParser) columnNameToIndex(str string) int {
	var l = len(str)
	if l < 1 || l > 2 {
		return -1
	}

	multi := l > 1
	// ascii not between 'A' - 'Z'  ( 65 - 90)
	b1 := int(str[0])
	b2 := 0
	if multi {
		b2 = int(str[1])
	}
	if b1 < 65 || b2 > 90 || b2 < 65 || b2 > 90 {
		return -1
	}

	if multi {
		return (b1-65+1)*26 + (b2 - 65 + 1)
	} else {
		return b1 - 65 + 1
	}
}

func (excel *ExcelParser) cellName(rowIndex int, colIndex int) string {
	return excel.columnIndexToName(colIndex) + strconv.Itoa(rowIndex)
}

func (excel *ExcelParser) pureText(obj interface{}) string {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.String:
		return obj.(string)
		// case reflect.Map:
		// 	_map := make(map[string]interface{})
		// 	if _map["richText"] {
		// 		arr := _map["richText"]
		// 		s := ""
		// 		if len(arr) > 0 {
		// 			for _, element : range arr {
		// 				if element["text"] {
		// 					if trimText {
		// 						s = s + strings.Trim(element["text"])
		// 					} else {
		// 						s = s + element["text"]
		// 					}
		// 				}
		// 			}
		// 			arr.forEach(func(element map[string]string) {
		// 			})
		// 		}
		// 		return s
		// 	}
	}

	panic("pureText")
}

func (excel *ExcelParser) columnIndex(key interface{}) int {
	switch reflect.TypeOf(key).Kind() {

	case reflect.Int:
		return key.(int)

	case reflect.String:
		str := key.(string)
		// column name to column index
		index := excel.columnNameToIndex(str)
		if index >= 0 {
			return index
		}

		return excel.TitleIndexes[str]

	default:
		panic("columnIndex")
		return -1
	}
}

// translate column name to uniformed excel column name like 'AZ'
// accepts one of the following forms:
//
//	column index: 1 .. 702
//	column title: 'student no', 'teacher name'
//	column alias: 'field1', 'field2'
//	standard excel column name: 'A' - 'Z', 'AA' - ... - 'ZZ'
func (excel *ExcelParser) columnName(key int) string {
	index := excel.columnIndex(key)
	return excel.columnIndexToName(index)
}

// Build title aliases data from excel column titles
func (excel *ExcelParser) buildTitleAliases(titles []string) {
	columnAliases := excel.Options.ColumnsAliases
	aliasMapping := make(AliasIndexes)

	excel.TitleAliases = AliasConfig{}
	excel.TitleIndexes = make(map[string]int)

	// build alias mapping
	if len(columnAliases) > 0 {
		for _, aliasGroup := range columnAliases {
			if len(aliasGroup) > 0 {
				excel.TitleAliases = append(excel.TitleAliases, aliasGroup)
				for _, alias := range aliasGroup {
					aliasMapping[alias] = aliasGroup
				}
			}
		}
	}

	// build titleAliases
	if len(titles) > 0 {
		index := excel.Options.InvalidColCount
		for _, title := range titles {
			excel.TitleIndexes[title] = index
			if len(title) > 0 && len(aliasMapping[title]) > 0 {
				for _, element := range aliasMapping[title] {
					excel.TitleIndexes[element] = index
				}
			}
			index++
		}
	}
}

func (excel *ExcelParser) titleAliasToColumnIndex(alias string) int {
	if len(excel.TitleAliases) > 0 {
		i := 1

		for _, aliases := range excel.TitleAliases {
			for _, val := range aliases {
				if val == alias {
					return i
				}
			}
			i++
		}
	}
	return -1
}

func (excel *ExcelParser) parseTitle(values []interface{}) []string {
	var titles = []string{}
	if excel.Options.CustomParseTitle != nil {
		f := excel.Options.CustomParseTitle
		titles = f(values)
	} else {
		for _, val := range values {
			titles = append(titles, val.(string))
		}
	}
	excel.buildTitleAliases(titles)
	return titles
}

func (excel *ExcelParser) loadData(sheetname string) (ExcelData, error) {
	excel.Data = ExcelData{}

	var file, err = xlsx.OpenFile(excel.Filepath)
	if err != nil {
		return excel.Data, err
	}
	firstCol := excel.Options.InvalidColCount

	var worksheet *xlsx.Sheet = nil
	if len(sheetname) > 0 {
		worksheet = file.Sheet[sheetname]
	} else {
		worksheet = file.Sheets[0]
	}
	if worksheet != nil {
		rows := worksheet.MaxRow
		cols := worksheet.MaxCol

		yellow(sheetname + ":" + strconv.Itoa(rows) + " x " + strconv.Itoa(cols))

		for i, row := range worksheet.Rows {
			rowNumber := i + 1

			// get row data
			var rowData = RowData{}
			for _, cell := range row.Cells {
				rowData = append(rowData, cell.Value)
			}

			// 数据不包含标题行和首列数据
			if rowNumber > excel.Options.TitleRowCount {
				excel.Data = append(excel.Data, rowData[firstCol:])
			} else {
				excel.TitleData = append(excel.TitleData, rowData[firstCol:])
				if rowNumber == excel.Options.TitleRowIndex {
					excel.ColumnTitles = excel.parseTitle(rowData[firstCol:])
				}
			}
		}

	}

	return excel.Data, nil
}

// check whether the xlsx contains specified sheet name
func (excel *ExcelParser) hasSheet(sheetname string) bool {

	var file, err = xlsx.OpenFile(excel.Filepath)
	if err != nil {
		printError(err)
		return false
	}

	worksheet := file.Sheet[sheetname]
	return worksheet != nil
}

func (excel *ExcelParser) getValue(row RowData, col string) string {
	_col := excel.columnIndex(col)
	if _col >= 0 && _col < len(row) {
		return stringlize(row[_col])
	}
	return ""
}
