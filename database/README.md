# goCommon database

## mysql
```
import (
    "github.com/obse4/goCommon/database"
    "time"
)

// 配置
var mysql = database.MysqlConfig{
	Name: "mysql测试",
	Username: "admin",
	Password: "123456",
	Database: "test",
	Url: "127.0.0.1",
	Port: "3306",
	Charset: "utf8mb4",
}
func main() {
    // 初始化
    database.InitMysqlConnect(&mysql)
    // 使用
    // mysql.Db 是gorm的db指针，用法参考gorm
    mysql.Db.Table("user").Create(&map[string]interface{}{
		"username": "obse4",
		"type": "root",
		"description": "测试用户",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
}

```

## postgres
```
import (
    "github.com/obse4/goCommon/database"
    "time"
)

// 配置
var postgres = database.PostgresConfig{
	Name:     "postgres测试",
	Username: "admin",
	Password: "123456",
	Database: "test",
	Url:      "127.0.0.1",
	Port:     "5432",
}
func main() {
    // 初始化
    database.InitPostgresConnect(&postgres)
    // 使用
    // postgres.Db 是gorm的db指针，用法参考gorm
    postgres.Db.Table("user").Create(&map[string]interface{}{
		"username": "obse4",
		"type": "root",
		"description": "测试用户",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
}
```
## redis
```
import (
    "github.com/obse4/goCommon/database"
    "time"
)

var redis = database.RedisConfig{
	Name: "redis测试",
	Url: "127.0.0.1",
	Port: "6379",
	Password: "",
}

func main()  {
	// 初始化
    database.InitRedisPool(&redis)
    // 使用
    // redis.Pool 是redigo的连接池指针，使用方法参考redigo
  	conn := database.GetRedisConn(redis.Pool, 0)
	conn.Do("RPUSH", "test_list", "测试内容")
	conn.Close()
}
```

## mongodb
```
import (
    "github.com/obse4/goCommon/database"
    "time"
)

var mongodb = database.MongoConfig{
	Name: "mongo测试",
	Url:  "127.0.0.1",
	Port: "27017",
}

type Order struct {
	Id         primitive.ObjectID `bson:"_id"`
	Price      int                `bson:"price"`
	OrderId    string             `bson:"order_id" mongo:"index"`
	Status     int                `bson:"status"`
	PayDate    string             `bson:"pay_date" mongo:"index; compound:'member_id_pay_date'"`
	RefundDate string             `bson:"refund_date" mongo:"compound:'member_id_refund_date'"`
	MemberId   string             `bson:"member_id" mongo:"compound:'member_id_pay_date','member_id_refund_date'"`
	CreatedAt  string             `bson:"created_at" mongo:"index"` 
	UpdatedAt  string             `bson:"updated_at"`              
}

func main() {
	// 初始化
	database.InitMongoDBConnect(&mongodb)
	// 注册表和索引
	database.AutoRegisterMongo(mongodb.Db, Order{})
	// 使用 mongodb.Db 是go.mongodb.org/mongo-driver/mongo驱动的数据库指针，使用方法参考mongo-driver
    // 写法1
	mongodb.Db.Collection("orders").InsertOne(context.TODO(), &bson.M{
		"price":       100,
		"order_id":    "20231222134530-0518",
		"status":      1,
		"pay_date":    "2023-12-22 13:45:30",
		"refund_date": "",
		"member_id":   "test0518",
		"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
	})
    // 写法2
    mongodb.Db.Collection(database.GetMongoCollection(mongodb.Db, Order{})).InsertOne(context.TODO(), &bson.M{
		"price":       100,
		"order_id":    "20231222134530-0518",
		"status":      1,
		"pay_date":    "2023-12-22 13:45:30",
		"refund_date": "",
		"member_id":   "test0518",
		"created_at":  time.Now().Format("2006-01-02 15:04:05"),
		"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
	})
}
```