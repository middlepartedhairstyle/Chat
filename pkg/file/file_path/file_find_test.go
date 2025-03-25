package file_path

import (
	"testing"
)

func TestGetWorkPath(t *testing.T) {
	path := GetWorkPath()
	t.Log(path)
}

func TestIsNull(t *testing.T) {
	b := IsNull("E:/golang/go/HiWe/config", "config.yaml")
	if !b {
		t.Errorf("err")
	}
}
