package utils

import (
	"strconv"
	"time"
)

// GetTimeUnixNanoS 返回string类型的时间戳
func GetTimeUnixNanoS() string {
	tm := time.Now().UnixNano()
	return strconv.FormatInt(tm, 10)
}

// GetTimeUnixNanoI 返回int类型的时间戳
func GetTimeUnixNanoI() int64 {
	tm := time.Now().UnixNano()
	return tm
}
