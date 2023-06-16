package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type UploadResult struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Students bool   `json:"students"`
	Schools  bool   `json:"schools"`
}

//	{
//	    "data": [
//	        {
//	            "_id": "a1",
//	            "key": 1,
//	            "fileName": "5fc1372e68a1010210d11e12635492e5c5.png",
//	            "filePath": "uploads/5fc1372e68a1010210d11e12635492e5c5.png",
//	            "uploader": "unknown",
//	            "timeOfUpload": "2023-04-21T06:55:45.493Z",
//	            "fileSize": 1
//	        }
//	    ]
//	}
type FileItem struct {
	Id           string `json:"id"`
	Key          int    `json:"key"`
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	FileSize     int64  `json:"fileSize"`
	TimeOfUpload string `json:"timeOfUpload"`
	Uploader     string `json:"uploader"`
}

type DeleteFileRequest struct {
	FilePath string `json:"filePath"`
	Id       string `json:"id"`
}

// Upload file
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// file from uploaded Multipart form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	srcFilename := header.Filename
	targetPath := "./Upload/" + srcFilename
	targetPath, _ = getUniqueFileName(targetPath)
	targetFilename := path.Base(targetPath)

	// Create uploaded file in disk
	out, err := os.Create(targetPath)
	if err != nil {
		httpResponseError(w, err)
		return
	}
	defer out.Close()

	// 复制文件数据
	fileSize, err := io.Copy(out, file)
	if err != nil {
		httpResponseError(w, err)
		return
	}

	fileUploaded := UploadResult{
		Name:     srcFilename,
		Filename: targetFilename,
		Size:     fileSize,
		Path:     targetPath,
	}

	//uploadedSpecialProcess(&fileUploaded)

	httpResponseJson(w, fileUploaded)
}

// list all files in folder
func listFileHandler(w http.ResponseWriter, r *http.Request) {
	fileList := []FileItem{}

	files, _ := enumerateFiles("Upload/", []string{".*"})
	for i, file := range files {
		fileItem := FileItem{
			Id:           strconv.Itoa(i + 1),
			Key:          i + 1,
			FileName:     filepath.Base(file),
			FilePath:     file,
			FileSize:     0,
			TimeOfUpload: getFileModifyTime(file),
			Uploader:     "",
		}
		fileList = append(fileList, fileItem)
	}

	httpResponseObject(w, fileList)
}

// delete file by path
func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	var request DeleteFileRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		httpResponseError(w, err)
		return
	}

	// remove file
	os.Remove(request.FilePath)

	httpResponseObject(w, "")
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	httpResponseFile(w, filePath)
}
