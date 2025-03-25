package jwt_rsa

import (
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/pkg/token/encrypt"
	"testing"
	"time"
)

func TestNewJwtRsa(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key"
	pr, pu := "pr.pem", "pu.pem"
	k, err := encrypt.ReadRSAKeyFile(path, pr, pu)
	if err != nil {
		t.Errorf("LoadRSAKeyFile error")
		return
	}
	jwtRsa := NewJwtRsa(map[string]interface{}{"time": time.Now().UnixNano(), "name": "hello"}, k)
	fmt.Println(jwtRsa)
	fmt.Println(jwtRsa.JwtRsaToToken())
	j := TokenToJwtRsa(jwtRsa.JwtRsaToToken())
	if j.Payload != jwtRsa.Payload || j.Signature != jwtRsa.Signature {
		t.Errorf("TokenToJwtRsa error")
	}

	if !j.JwtRsaVerify(k) {
		t.Errorf("TokenToJwtRsaVerify error")
	}
}

func TestCreateRSAKey(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key/"
	pr, pu := "pr.pem", "pu.pem"
	_, err := CreateRSAKey(512, path, pr, pu)
	if err != nil {
		t.Errorf("CreateRSAKey error")
		return
	}
}

func TestLoadRSAKeyFile(t *testing.T) {
	path := "E:/golang/go/HiWe/config/key"
	pr, pu := "pr.pem", "pu.pem"
	k := LoadRSAKeyFile(path, pr, pu)
	if k == nil {
		t.Errorf("LoadRSAKeyFile error")
	}
	fmt.Println(k)
}

func TestJwtRsa_DecodePayload(t *testing.T) {
	jwt := &JwtRsa{
		Payload:   "eyJpZCI6MTIzLCJuYW1lIjoiaGVsbG8iLCJ0eXBlIjoicmVmcmVzaF90b2tlbiJ9",
		Signature: "MDP44tdgnYw0w7ySmW1fiZnJidiCMreUzAyQP5x__CAd_-7XlTAqNiYoYZUXBYwnhReSjoDOEwHkea3Ea8Q1EQ",
	}
	fmt.Println(jwt.DecodePayloadToMap())
}

func TestJwtRsa_JwtRsaVerify(t *testing.T) {
	jwt := &JwtRsa{
		Payload:   "eyJpZCI6MTIzLCJuYW1lIjoiaGVsbG8iLCJ0eXBlIjoicmVmcmVzaF90b2tlbiJ9",
		Signature: "MDP44tdgnYw0w7ySmW1fiZnJidiCMreUzAyQP5x__CAd_-7XlTAqNiYoYZUXBYwnhReSjoDOEwHkea3Ea8Q1EQ",
	}
	path := "E:/golang/go/HiWe/config/key"
	pr, pu := "pr.pem", "pu.pem"
	k := LoadRSAKeyFile(path, pr, pu)
	if k == nil {
		t.Errorf("LoadRSAKeyFile error")
	}
	b := jwt.JwtRsaVerify(k)
	if !b {
		t.Errorf("JwtRsaVerify error")
	}
}
