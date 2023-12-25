package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtConfig struct {
	// secret
	SecretKey string
	// 过期时长 s
	ExpiresTime int64
}

type Claims interface {
	Valid() error
}

type StandardClaims struct {
	jwt.StandardClaims
}

// 创建jwt token
// claims为nil时默认使用对象StandardClaims
// 可自定义claims，需要继承StandardClaims
func (j *JwtConfig) CreateJwtToken(claims Claims) (jwtToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	if claims == nil {
		claims = StandardClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(j.ExpiresTime) * time.Second).Unix(),
			},
		}
		token.Claims = claims
	} else {
		token.Claims = claims
	}

	// !传参虽然定义了interface{}类型，但是如果不是[]byte会报错 `key is of invalid type`
	if j.SecretKey == "" {
		jwtToken, err = token.SigningString()
	} else {
		jwtToken, err = token.SignedString([]byte(j.SecretKey))
	}
	if err != nil {
		return "", err
	}
	return
}

// 解析token
// tokenString 字符串token
// claims 自定义继承StandardClaims的struct实例指针
func (j *JwtConfig) ParseJwtToken(tokenString string, claims Claims) error {
	var err error
	if tokenString == "" {
		return fmt.Errorf("无效的token")
	}

	// 解析token
	_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// !传参虽然定义了interface{}类型，但是如果不是[]byte会报错 `key is of invalid type`
		if j.SecretKey == "" {
			return nil, nil
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return err
	}
	return err
}
