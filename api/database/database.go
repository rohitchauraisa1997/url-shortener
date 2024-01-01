package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	fmt.Println("DB_ADDR", os.Getenv("DB_ADDR"))
	fmt.Println("DB_PASS", os.Getenv("DB_PASS"))
	fmt.Println("DB_USER", os.Getenv("DB_USER"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		Username: os.Getenv("DB_USER"),
		DB:       dbNo,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
