# goCommon config

```
import "github.com/obse4/goCommon/config"

func main() {
    // 配置文件夹及目录参考测试test文件夹
    type OtherConfig struct {
		IsTrue bool
		Type   string
	}
	var GlobalConfig struct {
		Name  string
		Id    int
		Other OtherConfig
		List  []OtherConfig
	}

    // config为文件夹名称，&GlobalConfig为配置结构体指针
    // CONFIG_MODE环境变量不配置默认使用env.yml配置文件
	config.InitConfig("config", &GlobalConfig)
}
```