package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type CommonResponse struct {
	Data interface{} `json:"data"`
}

func httpResponseError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func httpResponseFail(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpResponseJson(w http.ResponseWriter, obj interface{}) {
	jsonData, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func httpResponseObject(w http.ResponseWriter, obj interface{}) {
	resp := new(CommonResponse)
	resp.Data = obj
	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func httpResponseFile(w http.ResponseWriter, path string) {
	// 打开文件
	f, err := os.Open(path)
	if err != nil {
		httpResponseError(w, err)
		return
	}
	defer f.Close()

	// 设置 HTTP 响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(getFileSize(f)))

	// 将文件内容写入 HTTP 响应体
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
