package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(str string) string {
	return strings.ToUpper(Sha256(str))
}

func MakePasswordSha256(password string, salt string) string {
	return Sha256(password + salt)
}

func CheckPasswordSha256(password string, salt string, encryptionPassword string) bool {
	return Sha256(password+salt) == encryptionPassword
}

func MakeToken(email string, password string, tm string) string {
	token1 := Sha256(email + password)
	return Sha256(token1 + tm)
}
