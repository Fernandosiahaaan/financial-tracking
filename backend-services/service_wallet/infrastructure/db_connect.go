package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	connPostgrees "service-wallet/internal/store/postgree"
	connRedis "service-wallet/internal/store/redis"
	"service-wallet/utils"
	"time"

	"github.com/go-redis/redis"
)

func NewConnectionRedis(ctx context.Context) (*connRedis.RedisCln, error) {
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

	var redis *connRedis.RedisCln = &connRedis.RedisCln{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}

	// Ping the Redis server to check the connection
	_, err := client.Ping().Result()
	if err != nil {
		errMsg := utils.MessageError("infrastructure::NewConnectionRedis", fmt.Errorf("failed ping redis. err : %v", err))
		return nil, errMsg
	}

	return redis, nil
}

func NewConnectionPostgree(ctx context.Context) (*connPostgrees.StoreWalletPq, error) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		errMsg := utils.MessageError("infrastructure::NewConnectionPostgree", fmt.Errorf("could not connect to the database. err : %v", err))
		return nil, errMsg
	}

	dbCtx, dbCancel := context.WithCancel(ctx)
	return &connPostgrees.StoreWalletPq{
		Db:     db,
		Ctx:    dbCtx,
		Cancel: dbCancel,
	}, nil
}
