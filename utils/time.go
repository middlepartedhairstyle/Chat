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

// GetTimeToUTC 获取本地时间，将本地时间改为世界时区时间，默认为东8区时间
func GetTimeToUTC(utc ...time.Duration) time.Time {
	var u time.Duration = 8 //默认为东8区时间
	if len(utc) > 0 {
		u = utc[0]
	}
	now := time.Now().UTC()
	nowUTC := now.Add(u * time.Hour)
	return nowUTC
}
