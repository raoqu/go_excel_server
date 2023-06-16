package main

import (
	"io/fs"
	"net/http"
	"os"
)

func initEnvironment() {
	if _, err := os.Stat("./Upload"); !os.IsNotExist(err) {
		os.Mkdir("./Upload", fs.ModeDir)
	}
}

func main() {
	initEnvironment()

	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/files", listFileHandler)
	http.HandleFunc("/api/delete", deleteFileHandler)
	http.HandleFunc("/api/download", downloadFileHandler)
	http.HandleFunc("/api/excel_list", getExcelList)
	http.HandleFunc("/api/excel_data", getExcelData)
	http.ListenAndServe("127.0.0.1:8083", nil)
}
