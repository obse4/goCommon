# goCommon logger

## 初始化
```
import "github.com/obse4/goCommon/logger"

func main() {
    logger.InitLogger(&logger.LogConfig{
		LogOut:   true,
		StayDay:  1,
	})
}

```

## 使用
```
import "github.com/obse4/goCommon/logger"

func main() {
    logger.Debug("Hello %s", "world")
}
```

## 配置文件
查看配置`logger.LogConfig`