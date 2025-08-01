package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"service-user/internal/model"
	"service-user/utils"
	"time"

	"github.com/go-redis/redis"
)

const (
	PrefixKeyUserInfo  = "user-service:user"
	PrefixKeyLoginInfo = "user-service:jwt"
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

	var redis *RedisCln = &RedisCln{
		Redis:  client,
		Ctx:    ctxRedis,
		Cancel: cancelRedis,
	}

	// Ping the Redis server to check the connection
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("failed ping redis. err : %v", err)
	}

	return redis, nil
}

func (r *RedisCln) SaveUserInfo(user model.User) error {
	var msgErr error = nil
	userJson, err := json.Marshal(user)
	if err != nil {
		msgErr = utils.MessageError("Redis::SaveUserInfo", fmt.Errorf("failed convert user '%s' info to json", user.ID))
		return msgErr
	}

	// send user info to reddis data
	keyUserInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, user.ID)
	err = r.Redis.Set(keyUserInfo, userJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		msgErr = utils.MessageError("Redis::SaveUserInfo", fmt.Errorf("error saving login info to redis. err = %s", err.Error()))
		return msgErr
	}

	return nil
}

func (r *RedisCln) GetUserInfo(userId string) (user *model.User, err error) {
	var msgErr error = nil

	userInfo := fmt.Sprintf("%s:%s", PrefixKeyUserInfo, userId)
	userJson, err := r.Redis.Get(userInfo).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetUserInfo", fmt.Errorf("failed get user info from redis. err : %v", err))
		return nil, msgErr
	}
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		msgErr = utils.MessageError("Redis::GetUserInfo", fmt.Errorf("failed convert data user info from json. err : %v", err))
		return nil, msgErr
	}
	return user, nil
}

func (r *RedisCln) GetAllUserInfo() ([]model.User, error) {
	var msgErr error = nil
	keysPattern := fmt.Sprintf("%s:*", PrefixKeyUserInfo)

	keys, err := r.Redis.Keys(keysPattern).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetAllUserInfo", fmt.Errorf("failed to get keys from redis with pattern %s. err: %v", keysPattern, err))
		return nil, msgErr
	}

	// MGET all of data in keys
	usersJson, err := r.Redis.MGet(keys...).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetAllUserInfo", fmt.Errorf("failed to get multiple employments from redis. err: %v", err))
		return nil, msgErr
	}

	var users []model.User
	for _, userJson := range usersJson {
		if userJson == nil {
			continue
		}
		var user model.User
		err := json.Unmarshal([]byte(userJson.(string)), &user)
		if err != nil {
			msgErr = utils.MessageError("Redis::GetAllUserInfo", fmt.Errorf("failed to unmarshal employment json: %v", err))
			return nil, msgErr
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *RedisCln) GetLoginInfo(jwtToken string) (loginInfo model.LoginCacheData, err error) {
	var msgErr error = nil

	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	loginJson, err := r.Redis.Get(keyLoginInfo).Result()
	if err != nil {
		msgErr = utils.MessageError("Redis::GetLoginInfo", fmt.Errorf("failed get login info from redis. err : %v", err))
		return loginInfo, msgErr
	}
	err = json.Unmarshal([]byte(loginJson), &loginInfo)
	if err != nil {
		msgErr = utils.MessageError("Redis::GetLoginInfo", fmt.Errorf("failed convert data login info from json. err : %v", err))
		return loginInfo, msgErr
	}
	return loginInfo, nil
}

func (r *RedisCln) SetLoginInfo(ctx context.Context, jwtToken string, loginInfo model.LoginCacheData) error {
	var msgErr error = nil

	loginJson, err := json.Marshal(loginInfo)
	if err != nil {
		msgErr = utils.MessageError("Redis::GetLoginInfo", fmt.Errorf("failed convert login info '%s' to json. err : %v", loginInfo.Username, err))
		return msgErr
	}

	// send login info to reddis data
	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	err = r.Redis.Set(keyLoginInfo, loginJson, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
	if err != nil {
		msgErr = utils.MessageError("Redis::GetLoginInfo", fmt.Errorf("error saving login info '%s' to redis. err = %s", loginInfo.Username, err.Error()))
		return msgErr
	}
	return nil
}

func (r *RedisCln) DeleteLoginInfo(jwtToken string) error {
	var msgErr error = nil

	keyLoginInfo := fmt.Sprintf("%s:%s", PrefixKeyLoginInfo, jwtToken)
	err := r.Redis.Del(keyLoginInfo).Err()
	if err != nil {
		msgErr = utils.MessageError("Redis::DeleteLoginInfo", fmt.Errorf("error delete login info '%s' to redis. err = %s", jwtToken, err.Error()))
		return msgErr
	}

	return nil
}

func (r *RedisCln) Close() {
	r.Redis.Close()
	r.Cancel()
}
