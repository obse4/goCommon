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
    // 可配置环境变量CONFIG_PATH来控制配置文件路径，权重最高 // 例如/config/global.yml CONFIG_PATH=/config/global.yml
	// InitConfig函数配置的path权重其次
	// path为空，默认使用当前执行文件所在目录下的config.yml文件
	config.InitConfig("", &GlobalConfig)
}
```