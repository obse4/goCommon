# goCommon jwt

## 使用
```
import "github.com/obse4/goCommon/jwt"

func main() {
    // 初始化实例
    var jwtConfig = jwt.JwtConfig{
		ExpiresTime: 60,
		SecretKey:   "0518",
	}

    // 创建token
	token, err := jwtConfig.CreateJwtToken(nil)
    // 
    var res jwt.StandardClaims
    err := jwtConfig.ParseToken(token, &res)
    if err != nil {
		fmt.Printf("parse jwt token err %s", err.Error())
		return
	}
    fmt.Printf("parse token success:\n%#v\n", res)
}
```