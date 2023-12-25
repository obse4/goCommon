package util

import (
	"crypto/md5"
	"fmt"
)

func MD5(s string) string {
	// md5加密
	srcCode := md5.Sum([]byte(s))
	// 转换成16进制
	code := fmt.Sprintf("%x", srcCode)

	return string(code)
}
