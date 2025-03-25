package file

import (
	"github.com/middlepartedhairstyle/HiWe/pkg/file/file_path"
	"os"
	"path/filepath"
)

// CreateFile
// 创建文件,path为创建的文件路径,fileName为文件名称
func CreateFile(path string, fileName string) *os.File {
	//判断是否存在该文件
	//如果存在该文件
	if !file_path.IsNull(path, fileName) {
		return nil
	}
	//如果不存在该文件
	//构造完整文件路径
	filePath := filepath.Join(path, fileName)
	//创建文件
	file, _ := os.Create(filePath)
	return file
}

// CreateDir
// 构建文件夹,path为文件路径,dirName为构建的文件夹名称
func CreateDir(path string, dirName string) error {
	dirPath := filepath.Join(path, dirName)
	err := os.MkdirAll(dirPath, os.ModePerm)
	return err
}
