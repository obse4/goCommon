package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/obse4/goCommon/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Name     string          // 自定义名称
	Username string          // 用户名
	Password string          // 密码
	Url      string          // url链接
	Port     string          // 端口
	Database string          // 数据库名称
	Db       *mongo.Database // 数据库指针
}

// 初始化数据库
func InitMongoDBConnect(mongoose *MongoConfig) {
	var url string
	if mongoose.Username != "" {
		url = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoose.Username, mongoose.Password, mongoose.Url, mongoose.Port, mongoose.Database)
		if mongoose.Username == "admin" {
			url = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoose.Username, mongoose.Password, mongoose.Url, mongoose.Port)
		}
	} else {
		url = fmt.Sprintf("mongodb://%s:%s/%s", mongoose.Url, mongoose.Port, mongoose.Database)
	}

	if strings.Contains(mongoose.Url, "mongodb://") {
		// 支持仅填写url
		url = mongoose.Url
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		logger.Fatal("Mongo %s %s", mongoose.Name, err.Error())
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal("Mongo %s connect error %s", mongoose.Name, err.Error())
	}

	mongoose.Db = client.Database(mongoose.Database)
	logger.Info("Mongo %s 连接成功", mongoose.Name)
}

// 注册表和索引
// 自动注册的表名为struct的小写复数形式
// 标签 mongo
// 索引 index
// 复合索引 compound
// 示例
//
//	type Order struct {
//		Id          primitive.ObjectID `bson:"_id"`
//		Price       int                `bson:"price"`
//		OrderId     string             `bson:"order_id" mongo:"index"`
//		Status      int                `bson:"status"`
//		PayDate     string             `bson:"pay_date" mongo:"index; compound:'member_id_pay_date'"`
//		RefundDate  string             `bson:"refund_date" mongo:"compound:'member_id_refund_date'"`
//		MemberId    string             `bson:"member_id" mongo:"compound:'member_id_pay_date','member_id_refund_date'"`
//		CreatedAt   string             `bson:"created_at" mongo:"index"` // 添加时间 YYYY-MM-dd hh:mm:ss
//		UpdatedAt   string             `bson:"updated_at"`               // 更新时间 YYYY-MM-dd hh:mm:ss
//	}
func AutoRegisterMongo(db *mongo.Database, v ...any) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, val := range v {
		t := reflect.TypeOf(val)

		p := pluralize.NewClient()
		plural := p.Plural(t.Name())
		colName := strings.ToLower(plural[:1] + plural[1:])
		indexList := []string{}
		compoundMap := make(map[string]map[string]any)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			f := field.Tag.Get("bson")
			if f == "" {
				logger.Fatal("Mongo", "字段标签bson不能为空")
			}
			mctx := field.Tag.Get("mongo")
			if mctx != "" {
				// 处理mongo定义索引
				// 切分
				slice := []string{}
				for _, v := range strings.Split(mctx, ";") {
					slice = append(slice, strings.TrimSpace(v))
				}

				for _, v := range slice {
					// 独立索引
					if v == "index" {
						indexList = append(indexList, f)
					}

					// 联合索引
					if strings.Contains(v, "compound") {

						comp := strings.ReplaceAll(v, "compound:", "")

						var compList []string

						for _, val := range strings.Split(comp, ",") {
							compList = append(compList, strings.ReplaceAll(strings.TrimSpace(val), "'", ""))
						}

						for _, compName := range compList {
							if compoundMap[compName] == nil {
								compoundMap[compName] = make(map[string]any)
							}
							compoundMap[compName][f] = nil
						}
					}
				}
			}
		}

		// 校验联合索引长度
		// 小于2则panic
		var indexModelList []mongo.IndexModel
		for _, compItem := range compoundMap {
			if len(compItem) < 2 {
				logger.Fatal("Mongo", "联合索引字段不足")
			}
			keys := bson.D{}
			for indexItem := range compItem {
				keys = append(keys, bson.E{Key: indexItem, Value: 1})
			}
			indexModelList = append(indexModelList, mongo.IndexModel{Keys: keys, Options: nil})
		}

		for _, v := range indexList {
			indexModelList = append(indexModelList, mongo.IndexModel{Keys: bson.D{{Key: v, Value: 1}}, Options: nil})

		}

		db.CreateCollection(ctx, colName)
		col := db.Collection(colName)
		indexes, err := col.Indexes().CreateMany(ctx, indexModelList)
		if err != nil {
			logger.Fatal("Mongo '%s' Index Create Fail: %s", colName, err.Error())
		}
		logger.Info("Mongo '%s' Index Create Success: %v", colName, indexes)
	}
}

func GetMongoCollection(db *mongo.Database, v interface{}) *mongo.Collection {
	t := reflect.TypeOf(v)
	p := pluralize.NewClient()
	plural := p.Plural(t.Name())
	name := strings.ToLower(plural[:1] + plural[1:])
	return db.Collection(name)
}
