package file

import "testing"

func TestCreateFile(t *testing.T) {
	path := "E:/golang/go/HiWe/config/"
	name := "key"
	CreateFile(path, name)
}

func TestCreateDir(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key/key"
	name := ""
	err := CreateDir(path, name)
	if err != nil {
		t.Errorf("err%s", err.Error())
	}
}
