# goCommon

服务端开发常用包封装，简化配置文件，便于开发者专注业务逻辑

## comm
常用包集成，参考配置文件`config_example.yml`修改和删除不需要的配置即可


```
import (
    comm "github.com/obse4/goCommon"
)

func main() {
    comm.Init("/github/goCommon/config.yml")
    comm.Run()
}
```

## logger
简单易配置的文件或控制台日志输出包

[使用文档](https://github.com/obse4/goCommon/blob/main/logger/README.md)

## database
常用数据库驱动及常用配置封装包

[使用文档](https://github.com/obse4/goCommon/blob/main/database/README.md)

## httpserver
gin驱动及常用中间件封装包

[使用文档](https://github.com/obse4/goCommon/blob/main/httpserver/README.md)

## config
配置文件github.com/spf13/viper的二次封装

[使用文档](https://github.com/obse4/goCommon/blob/main/config/README.md)

## util
其他常用包二次封装

[使用文档](https://github.com/obse4/goCommon/blob/main/util/README.md)

## kafka
`github.com/IBM/sarama`二次封装

[使用文档](https://github.com/obse4/goCommon/blob/main/kafka/README.md)