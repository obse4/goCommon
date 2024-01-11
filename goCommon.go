package comm

import (
	"github.com/gin-gonic/gin"
	"github.com/obse4/goCommon/config"
	"github.com/obse4/goCommon/database"
	"github.com/obse4/goCommon/httpserver"
	"github.com/obse4/goCommon/kafka"
	"github.com/obse4/goCommon/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	HttpServer    *gin.Engine
	Mysql         = make(map[string]*gorm.DB)
	Mongodb       = make(map[string]*mongo.Database)
	Postgres      = make(map[string]*gorm.DB)
	KafkaProducer = make(map[string]*kafka.Producer)
	KafkaConsumer = make(map[string]*kafka.Consumer)
)

type goCommonConfig struct {
	HttpServer    httpserver.HttpServerConfig `yaml:"httpServer"`
	Log           logger.LogConfig            `yaml:"log"`
	Mysql         []database.MysqlConfig      `yaml:"mysql"`
	Mongodb       []database.MongoConfig      `yaml:"mongodb"`
	Postgres      []database.PostgresConfig   `yaml:"postgres"`
	KafkaProducer []kafka.KafkaProducerConfig `yaml:"kafkaProducer"`
	KafkaConsumer []kafka.KafkaConsumerConfig `yaml:"kafkaConsumer"`
}

var initConfig goCommonConfig

// 初始化
// 建议使用CONFIG_PATH环境变量，权重最高
// configPath字段为空字符串，将使用执行文件所在目录下的config.yml作为执行文件
// configPath需要使用配置文件的绝对路径
func Init(configPath string) {
	// 初始化配置
	config.InitConfig(configPath, &initConfig)

	// 初始化日志
	logger.InitLogger(&initConfig.Log)

	// 初始化http服务
	HttpServer = httpserver.NewHttpServer(&initConfig.HttpServer)

	// 初始化数据库
	for _, v := range initConfig.Mysql {
		Mysql[v.Name] = database.InitMysqlConnect(&v)
	}
	for _, v := range initConfig.Mongodb {
		Mongodb[v.Name] = database.InitMongoDBConnect(&v)
	}
	for _, v := range initConfig.Postgres {
		Postgres[v.Name] = database.InitPostgresConnect(&v)
	}

	// 初始化kafka
	for _, v := range initConfig.KafkaProducer {
		KafkaProducer[v.Name] = kafka.NewKafkaProducer(&v)
	}
	for _, v := range initConfig.KafkaConsumer {
		KafkaConsumer[v.Name] = kafka.NewKafkaConsumer(&v)
	}
}

// 启动
func Run() {
	initConfig.HttpServer.Init()
}
