package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getUniqueFileName(name string) (string, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return name, nil
	}

	i := 1
	for {
		newName := name[:len(name)-len(filepath.Ext(name))] + "_" + strconv.Itoa(i) + filepath.Ext(name)
		if _, err := os.Stat(newName); os.IsNotExist(err) {
			return newName, nil
		}
		i++
	}
}

func getFileModifyTime(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return ""
	}
	createTime := fileInfo.ModTime()
	return createTime.Format(time.RFC822)
}

func pathMatchExts(path string, exts []string) bool {
	for _, ext := range exts {
		if ext == ".*" || filepath.Ext(path) == ext {
			return true
		}
	}
	return false
}

// enumerate files
//
//	exts: with "." prefix, e.g. {".xls", ".xlsx"}
func enumerateFiles(folderPath string, exts []string) ([]string, error) {
	var xlsxFiles = []string{}

	var err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		// Check if required files
		base := filepath.Base(path)
		if info == nil || info.IsDir() || strings.HasPrefix(base, ".") || base == "__merged.xlsx" {
			return nil
		} else if pathMatchExts(path, exts) {
			xlsxFiles = append(xlsxFiles, path)
			return nil
		}
		return nil
	})
	return xlsxFiles, err
}

func removeFiles(filePaths []string) {
	for _, file := range filePaths {
		os.Remove(file)
	}
}

func fileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir() && fileInfo.Mode().IsRegular()
}

func copyFile(destPath string, srcPath string) {
	os.Remove(destPath)

	src, _ := os.Open(srcPath)
	dest, _ := os.Create(destPath)

	io.Copy(dest, src)

	defer src.Close()
	defer dest.Close()
}

func getFileSize(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}
