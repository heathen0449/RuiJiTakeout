package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Context context.Context

func RedisInit() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func RedisClose() {
	err := Rdb.Close()
	if err != nil {
		panic(err)
	}
}
