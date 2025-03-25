package file_path

import (
	"os"
	"path/filepath"
)

//GetWorkPath
/*
该函数用于获取本程序运行的工作目录
*/
func GetWorkPath() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return currentDir
}

//IsNull
/*
判断指定目录下是否存在文件或者文件夹
*/
func IsNull(path string, fileName string) bool {
	file := filepath.Join(path, fileName)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		return false
	}
	return false
}
