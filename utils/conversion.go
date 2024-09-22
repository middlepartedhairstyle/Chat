package utils

import "strconv"

// ToUint64 将string转换为Uint64
func ToUint64(str string) uint64 {
	u, _ := strconv.ParseUint(str, 10, 64)
	return u
}
