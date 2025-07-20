package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"service-user/infrastructure/redis"
	"service-user/internal/model"
	"service-user/internal/model/request"
	"service-user/internal/model/response"
	"service-user/internal/store"
	"service-user/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

var flagCacthAllUsers bool = false

type UserService struct {
	repo   *store.UserStore
	ctx    context.Context
	cancel context.CancelFunc
	redis  *redis.RedisCln
}

func NewUserService(ctx context.Context, redis *redis.RedisCln, repo *store.UserStore) *UserService {
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return &UserService{
		repo:   repo,
		ctx:    serviceCtx,
		cancel: serviceCancel,
		redis:  redis,
	}
}

func (s *UserService) CreateNewUser(createUser request.CreateUserRequest) (bodyResp response.ResponseHttp, err error) {
	var msgErr error = nil
	createUser.Password = strings.TrimSpace(createUser.Password)
	hashPassword, err := s.HashPassword(createUser.Password)
	if err != nil {
		msgErr = utils.MessageError("Service::CreateNewUser", fmt.Errorf("failed hash password. err : %v", err))
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed hash password with username '%s' [E001]", createUser.Username), MessageErr: msgErr.Error()}, msgErr
	}

	createUser.Password = hashPassword
	existUser, err := s.repo.GetUserByName(createUser.Username)
	if err != nil {
		msgErr = utils.MessageError("Repository::GetUserByName", err)
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed read user with username '%s' [E002]", createUser.Username), MessageErr: msgErr.Error()}, msgErr
	} else if existUser != nil {
		err = fmt.Errorf("user with username '%s' already created [E003]", existUser.Username)
		msgErr = utils.MessageError("Service::CreateNewUser", err)
		return response.ResponseHttp{IsError: true, Message: err.Error(), MessageErr: msgErr.Error()}, msgErr
	}

	var user model.User = model.User{
		ID:          uuid.New().String(),
		Username:    createUser.Username,
		Password:    createUser.Password,
		FullName:    createUser.FullName,
		Email:       createUser.Email,
		PhoneNumber: createUser.PhoneNumber,
		Role:        createUser.Role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	userID, err := s.repo.CreateNewUser(user)
	if err != nil {
		msgErr = utils.MessageError("Repository::CreateNewUser", err)
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed create user with username '%s' [E004]", createUser.Username), MessageErr: msgErr.Error()}, msgErr
	}

	if err = s.redis.SaveUserInfo(user); err != nil {
		msgErr = utils.MessageError("Redis::SaveUserInfo", err)
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed save user with username '%s' to redis [E005]", createUser.Username), MessageErr: msgErr.Error()}, msgErr
	}

	user.ID = userID
	return response.ResponseHttp{IsError: false, Message: fmt.Sprintf("success save user with username '%s' ", createUser.Username), Data: user}, nil
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	var usersInfo []model.User = []model.User{}
	var err error

	if flagCacthAllUsers {
		usersInfo, err = s.redis.GetAllUserInfo()
		if (err == nil) && (usersInfo != nil) {
			return usersInfo, nil
		}
	}

	usersInfo, err = s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range usersInfo {
		if err = s.redis.SaveUserInfo(user); err != nil {
			fmt.Printf("failed set to redis data education with name %s\n", user.Username)
			flagCacthAllUsers = false
		}
	}
	flagCacthAllUsers = true

	return usersInfo, nil
}

func (s *UserService) GetUserByName(user model.User) (model.User, error) {
	existUser, err := s.repo.GetUserByName(user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return *existUser, fmt.Errorf("username not found")
		}
		return *existUser, fmt.Errorf("error sql. err = %s", err.Error())
	}

	// Verifikasi apakah password cocok dengan hash
	match := s.VerifyPassword(user.Password, existUser.Password)
	if !match {
		return *existUser, fmt.Errorf("password not equal")
	}

	return *existUser, nil
}

func (s *UserService) GetUserById(userId string) (*model.User, error) {
	exitUser, err := s.redis.GetUserInfo(userId)
	if (err == nil) && (exitUser != nil) {
		return exitUser, nil
	}

	existUser, err := s.repo.GetUserById(userId)
	if err != nil && err != sql.ErrNoRows {
		return existUser, utils.MessageError("Repo::GetUserById", err)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return existUser, nil
}

func (s *UserService) UpdateUser(user model.User) (response.ResponseHttp, error) {
	var errMsg error = nil
	user.Password = strings.TrimSpace(user.Password)
	hashPassword, err := s.HashPassword(user.Password)
	if err != nil {
		errMsg = utils.MessageError("Repo::GetUserById", err)
		return response.ResponseHttp{IsError: true, Message: "failed hash password from system", MessageErr: errMsg.Error()}, errMsg
	}
	user.Password = hashPassword
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = model.RoleUser
	}

	id, err := s.repo.UpdateUser(user)
	if err != nil {
		errMsg = utils.MessageError("Repo::UpdateUser", err)
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed update user '%s' to db", user.Username), MessageErr: errMsg.Error()}, errMsg
	}
	user.ID = id

	if err = s.redis.SaveUserInfo(user); err != nil {
		errMsg = utils.MessageError("Repo::UpdateUser", err)
		return response.ResponseHttp{IsError: true, Message: fmt.Sprintf("failed update user '%s' to redis", user.Username), MessageErr: errMsg.Error()}, errMsg
	}
	return response.ResponseHttp{IsError: false, Message: fmt.Sprintf("success update user '%s'", user.Username), Data: user}, nil
}

func (s *UserService) Close() {
	s.cancel()
}
