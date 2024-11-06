package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

// StringToUint64 将string转换为Uint64
func StringToUint64(str string) uint64 {
	u, _ := strconv.ParseUint(str, 10, 64)
	return u
}

// StringToUint string转uint
func StringToUint(str string) (uint, error) {
	unsignedInt64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0, err
	}

	// 转换为 uint (适用于32位和64位平台)
	var unsignedInt uint = uint(unsignedInt64)
	return unsignedInt, nil
}

// StringToUint8 string转uint8
func StringToUint8(str string) (uint8, error) {
	unsignedInt64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return 0, err
	}

	// 转换为 uint (适用于32位和64位平台)
	var unsignedInt uint8 = uint8(unsignedInt64)
	return unsignedInt, nil
}

// UintToBytes uint转[]byte
func UintToBytes(num uint) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num) // 使用大端字节序
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
