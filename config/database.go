package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

func InitDataBase() {
	dsn := "user=arifu password=arifu dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Taipei"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// TODO : database error
		fmt.Println("db error")
	}

	DataBase = db
}

func GetDb() **gorm.DB {
	return &DataBase
}

// Redis connection

func RedisNewConnection() **redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &rdb
}
