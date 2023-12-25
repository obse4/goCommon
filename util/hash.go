package util

import (
	"golang.org/x/crypto/bcrypt"
)

func String2Hash(v string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost) //加密处理
	return string(hash)
}

func CompareHash(old, new string) (bool, error) {
	var same = false
	err := bcrypt.CompareHashAndPassword([]byte(old), []byte(new))

	if err == nil {
		same = true
	}

	return same, err
}
