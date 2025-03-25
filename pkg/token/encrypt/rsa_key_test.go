package encrypt

import (
	"fmt"
	"testing"
)

func TestNewRSAKey(t *testing.T) {
	k := NewRSAKey(512)
	if k == nil {
		t.Errorf("NewRSAKey error")
		return
	}
	if k.PrivateKey == nil || k.PublicKey == nil {
		t.Errorf("NewRSAKey error")
	}
}

func TestRSAKey_SaveRSAKeyFile(t *testing.T) {
	path := "E:/golang/go/HiWe/config/"
	pr, pu := "pr.pem", "pu.pem"
	k := NewRSAKey(512)
	err := k.SaveRSAKeyFile(path, pr, pu)
	if err != nil {
		fmt.Println(err)
		t.Errorf("SaveRSAKeyFile error")
	}
}

func TestRSAKey_LoadRSAKeyFile(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key/"
	pr, pu := "pr.pem", "pu.pem"
	k, err := ReadRSAKeyFile(path, pr, pu)
	if err != nil {
		t.Errorf("LoadRSAKeyFile error")
		return
	}
	fmt.Println(k)
}

func TestRSAKey_SignPKCS1v15(t *testing.T) {
	path := "E:/golang/go/HiWe/config/"
	pr, pu := "pr.pem", "pu.pem"
	str := "hello"
	k, err := ReadRSAKeyFile(path, pr, pu)
	if err != nil {
		t.Errorf("LoadRSAKeyFile error")
		return
	}
	v15, err := k.SignPKCS1v15(str)
	if err != nil {
		t.Errorf("SignPKCS1v15 error")
		fmt.Println(err)
		return
	}
	fmt.Println(v15)
}

func TestRSAKey_VerifyPKCS1v15(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key"
	pr, pu := "pr.pem", "pu.pem"
	str := "helloprivateKey, err := rsa.GenerateKey(rand.Reader, 2048)"
	k, err := ReadRSAKeyFile(path, pr, pu)
	if err != nil {
		t.Errorf("LoadRSAKeyFile error")
		return
	}
	v15, err := k.SignPKCS1v15(str)
	if err != nil {
		return
	}
	err = k.VerifyPKCS1v15(str, v15)
	if err != nil {
		t.Errorf("VerifyPKCS1v15 error")
		return
	}
	fmt.Println(v15)
}
