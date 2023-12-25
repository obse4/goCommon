# goCommon util

## 使用

- jwt

```
import "github.com/obse4/goCommon/util"

func main() {
    // 初始化实例
    var jwtConfig = util.JwtConfig{
		ExpiresTime: 60,
		SecretKey:   "0518",
	}

    // 创建token
	token, err := jwtConfig.CreateJwtToken(nil)
    // 
    var res util.StandardClaims
    err := jwtConfig.ParseToken(token, &res)
    if err != nil {
		fmt.Printf("parse jwt token err %s", err.Error())
		return
	}
    fmt.Printf("parse token success:\n%#v\n", res)
}
```

- aes

```
import "github.com/obse4/goCommon/util"

func main() {
    aes, err := util.NewAes("HELLOWORLDGOPHER")
    if err != nil {
		fmt.Printf("create aes struct err %s", err.Error())
		return 
	}
    aesEncode, err := aes.EnCode("^-^")

	if err != nil {
		fmt.Printf("encode aes err %s", err.Error())
		return 
	}

    fmt.Printf("aes encode %s\n", aesEncode)

   aesDecode, err := aes.Decode("eLVsUNgHirfIVv+qU5OmKQ==")

	if err != nil {
		fmt.Printf("decode aes err %s", err.Error())
		return
	}

	fmt.Printf("aes decode %s\n", aesDecode) 

}
```

- md5

```
import "github.com/obse4/goCommon/util"

func main() {
    res := util.MD5("^-^")
	fmt.Printf("md5 %s\n", res)
}
```
- hash
  
```
import "github.com/obse4/goCommon/util"

func main() {
    res := util.String2Hash("^-^")

	fmt.Printf("hash old %s\n", res)

	same, err := util.CompareHash(res, "^-^")

	if err != nil {
		fmt.Printf("hash compare err %s\n", err.Error())
		return
	}

	fmt.Printf("the hash compare is %v\n", same)
}
```