package db

import (
	"context"
	"log"
	"time"

	"github.com/NgeKaworu/to-do-list-go/src/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MongoClient mongo client
type MongoClient struct {
	MgEngine *mongo.Client //mongo client
	Mdb      string
}

// NewMongoClient 实例工厂
func NewMongoClient() *MongoClient {
	return &MongoClient{}
}

// Open 开启连接池
func (d *MongoClient) Open(mg, mdb string, initdb bool) error {
	d.Mdb = mdb
	ops := options.Client().ApplyURI(mg)
	p := uint64(39000)
	ops.MaxPoolSize = &p
	ops.WriteConcern = writeconcern.New(writeconcern.J(true), writeconcern.W(1))
	ops.ReadPreference = readpref.PrimaryPreferred()
	db, err := mongo.NewClient(ops)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Connect(ctx)
	if err != nil {
		return err
	}

	//err = db.Ping(ctx, readpref.PrimaryPreferred())
	//if err != nil {
	//	log.Println("ping err", err)
	//}

	d.MgEngine = db

	//初始化数据库
	if initdb {
		var session *mongo.Client
		session, err = mongo.NewClient(ops)
		if err != nil {
			panic(err)
		}
		err = session.Connect(context.Background())
		if err != nil {
			panic(err)
		}
		defer session.Disconnect(context.Background())

		// 记录表
		t := session.Database(mdb).Collection(models.TTask)
		indexView := t.Indexes()
		_, err := indexView.CreateMany(context.Background(), []mongo.IndexModel{
			{Keys: bsonx.Doc{bsonx.Elem{Key: "uid", Value: bsonx.Int32(1)}}},
			{Keys: bsonx.Doc{bsonx.Elem{Key: "createAt", Value: bsonx.Int32(-1)}}},
			{Keys: bsonx.Doc{bsonx.Elem{Key: "updateAt", Value: bsonx.Int32(-1)}}},
		})

		if err != nil {
			log.Println(err)
		}

	}

	return nil
}

// GetColl 获取表名
func (d *MongoClient) GetColl(coll string) *mongo.Collection {
	col, _ := d.MgEngine.Database(d.Mdb).Collection(coll).Clone()
	return col
}

// Close 关闭连接池
func (d *MongoClient) Close() {
	d.MgEngine.Disconnect(context.Background())
}
