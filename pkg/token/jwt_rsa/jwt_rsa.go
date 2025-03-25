package jwt_rsa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/middlepartedhairstyle/HiWe/pkg/token/encrypt"
	"os"
)

//JwtRsa
/*
JwtRsa为以Rsa加密实现的token格式
Payload为有效载荷字段，使用时应保证该字段的唯一性
Signature为签名字段，是对Payload字段的签名用于身份校验
*/
type JwtRsa struct {
	Payload   string
	Signature string
}

//NewJwtRsa
/*
实例化一个JwtRsa，需要传入payload数据(确保该数据的唯一性)，传入一对rsaKey密钥
*/
func NewJwtRsa(payload map[string]interface{}, rsaKey *encrypt.RSAKey) *JwtRsa {
	jwtRSA := &JwtRsa{}
	pl, _ := json.Marshal(payload)
	jwtRSA.Payload = base64.RawURLEncoding.EncodeToString(pl)
	jwtRSA.JwtRsaSign(rsaKey)
	return jwtRSA
}

//CreateRSAKey
/*
创建一对指定长度RSA密钥，并存储在指定路径
*/
func CreateRSAKey(bits int, path string, prFileName string, puFileName string) (*encrypt.RSAKey, error) {
	//创建一对RSA密钥
	key := encrypt.NewRSAKey(bits)
	//创建路径
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	//将密钥存储进文件
	err = key.SaveRSAKeyFile(path, prFileName, puFileName)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func LoadRSAKeyFile(path string, prFileName string, puFileName string) *encrypt.RSAKey {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	key, err := encrypt.ReadRSAKeyFile(path, prFileName, puFileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return key
}

func JwtRsaGetPayload(token string) string {
	payload, n := "", len(token)
	for i := n - 1; i >= 0; i-- {
		if token[i] == '.' {
			payload = token[:i]
			return payload
		}
	}
	return payload
}

func JwtRsaGetSignature(token string) string {
	signature, n := "", len(token)
	for i := n - 1; i >= 0; i-- {
		if token[i] == '.' {
			signature = token[i+1:]
			return signature
		}
	}
	return signature
}

func TokenToJwtRsa(token string) *JwtRsa {
	jwtRSA := &JwtRsa{}
	jwtRSA.Payload = JwtRsaGetPayload(token)
	jwtRSA.Signature = JwtRsaGetSignature(token)
	return jwtRSA
}

func (jwtRSA *JwtRsa) DecodePayloadToMap() map[string]interface{} {
	decode, _ := base64.RawURLEncoding.DecodeString(jwtRSA.Payload)
	var data map[string]interface{}
	err := json.Unmarshal(decode, &data)
	if err != nil {
		return nil
	}
	return data
}

func (jwtRSA *JwtRsa) JwtRsaSign(rsaKey *encrypt.RSAKey) {
	jwtRSA.Signature, _ = rsaKey.SignPKCS1v15(jwtRSA.Payload)
}

func (jwtRSA *JwtRsa) JwtRsaVerify(rsaKey *encrypt.RSAKey) bool {
	if rsaKey.VerifyPKCS1v15(jwtRSA.Payload, jwtRSA.Signature) != nil {
		return false
	}
	return true
}

func (jwtRSA *JwtRsa) JwtRsaToToken() string {
	return jwtRSA.Payload + "." + jwtRSA.Signature
}
