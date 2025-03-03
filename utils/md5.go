package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5 md5加密返回小写字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// MD5 md5加密返回大写字符串
func MD5(srt string) string {
	return strings.ToUpper(Md5(srt))
}

func MakePasswordMd5(password string, salt string) string {
	return Md5(password + salt)
}

func CheckPasswordMd5(password string, salt string, encryptionPassword string) bool {
	return Md5(password+salt) == encryptionPassword
}

// MakeTokenMd5 产生token,利用username和password产生token1，再利用token1和sale产生token2
func MakeTokenMd5(email string, password string, salt string) string {
	var token = MD5(email + password)
	return Md5(token + salt)
}
