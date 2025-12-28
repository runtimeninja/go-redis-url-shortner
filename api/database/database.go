package database

import (
	"context"
	"os"
	"time"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	addr := os.Getenv("DB_ADDR")
	if addr == "" {
		panic("DB_ADDR not set in environment")
	}

	pass := os.Getenv("DB_PASS")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbNo,
	})

	// ensure redis is reachable
	timeoutCtx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(timeoutCtx).Err(); err != nil {
		panic("redis connection failed: " + err.Error())
	}

	return client
}
