package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"service-user/internal/model"
	"time"

	"github.com/go-redis/redis"
)

const (
	PrefixKeyUserInfo = "user-service:user"
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
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	var redis *RedisCln = &RedisCln{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}
	return redis, nil
}

func (r *RedisCln) SaveUserInfo(user model.User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed convert user info to json")
	}

	// send user info to reddis data
	keyUserInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, user.Id)
	err = r.Redis.Set(keyUserInfo, userJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		return fmt.Errorf("error saving login info to redis. err = %s", err.Error())
	}
	return nil
}

func (r *RedisCln) GetUserInfo(userId string) (user *model.User, err error) {
	userInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, userId)
	userJson, err := r.Redis.Get(userInfo).Result()
	if err != nil {
		return nil, fmt.Errorf("failed get user info from redis")
	}
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return nil, fmt.Errorf("failed convert data user info from json")
	}
	return user, nil
}

func (r *RedisCln) Close() {
	r.Redis.Close()
	r.Cancel()
}
