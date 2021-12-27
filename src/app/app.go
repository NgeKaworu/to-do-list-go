package app

import (
	"github.com/NgeKaworu/to-do-list-go/src/db"
	"github.com/go-redis/redis/v8"
)

// App
type App struct {
	uc    *string
	mongo *db.MongoClient
	rdb   *redis.Client
}

// New 工厂方法
func New(
	uc *string,
	mongo *db.MongoClient,
	rdb *redis.Client,
) *App {

	return &App{
		uc,
		mongo,
		rdb,
	}
}
