package utils

import (
	"os"
	"path/filepath"
)

func CreateFile(path string, name string) (*os.File, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(exePath)
	fileDir := filepath.Join(dir, path)
	if _, err = os.Stat(fileDir); os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(filepath.Join(fileDir, name))
	return file, err
}

func CreateFilePath(path string) (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exePath)
	fileDir := filepath.Join(dir, path)
	if _, err = os.Stat(fileDir); os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return fileDir, err
}
