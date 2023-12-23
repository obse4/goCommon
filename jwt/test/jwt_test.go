package test

import (
	"fmt"
	"testing"

	"github.com/obse4/goCommon/jwt"
)

func TestDefaultCreateJwtToken(t *testing.T) {
	var jwtConfig = jwt.JwtConfig{
		ExpiresTime: 60,
		SecretKey:   "0518",
	}

	token, err := jwtConfig.CreateJwtToken(nil)

	if err != nil {
		fmt.Printf("create jwt token err %s", err.Error())
		t.Fail()
		return
	}

	fmt.Printf("create token success:\n%s\n", token)
}

func TestDefaultParseJwtToken(t *testing.T) {
	var jwtConfig = jwt.JwtConfig{
		ExpiresTime: 60,
		SecretKey:   "0518",
	}

	var res jwt.StandardClaims

	err := jwtConfig.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDMzMTc0ODd9.-RH3bk9DJbk1zcz_rSvp1jQdfxYOeGdU1oE9bLinj80", &res)

	if err != nil {
		fmt.Printf("parse jwt token err %s", err.Error())
		t.Fail()
		return
	}

	fmt.Printf("parse token success:\n%#v\n", res)
}
