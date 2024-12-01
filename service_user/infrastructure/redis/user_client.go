package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type RedisCln struct {
	Redis  *redis.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewReddisClient(ctx context.Context) (*RedisCln, error) {
	// Connect to Redis
	ctxRedis, cancelRedis := context.WithCancel(ctx)
	host := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	var opts *redis.Options = &redis.Options{
		Addr:        host, // Replace with your Redis server address
		Password:    "",   // No password for local development
		DB:          0,    // Default DB
		DialTimeout: 10 * time.Second,
	}
	client := redis.NewClient(opts)

	// Ping the Redis server to check the connection
	pong, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	var redis *RedisCln = &RedisCln{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}
	fmt.Println("Connected to Redis:", pong)
	return redis, nil
}

func (r *RedisCln) Close() {
	r.Redis.Close()
	r.Cancel()
}
