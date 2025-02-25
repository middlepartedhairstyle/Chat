package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Sha256 返回小写的sha256
func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 返回大写的SHA256
func SHA256(str string) string {
	return strings.ToUpper(Sha256(str))
}

// MakePasswordSha256 加密
func MakePasswordSha256(password string, salt string) string {
	return Sha256(password + salt)
}

// CheckPasswordSha256 校验密码正确性
func CheckPasswordSha256(password string, salt string, encryptionPassword string) bool {
	return Sha256(password+salt) == encryptionPassword
}

// MakeToken 返回token
func MakeToken(email string, password string, tm string) string {
	token1 := Sha256(email + password)
	return Sha256(token1 + tm)
}
