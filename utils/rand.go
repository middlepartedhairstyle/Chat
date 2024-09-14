package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// RandNum 基础随机数
func RandNum() uint32 {
	stringCode := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	code, _ := strconv.Atoi(stringCode)
	return uint32(code)
}

// RandString 基础随机数，字符串类型,6位
func RandString() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}
