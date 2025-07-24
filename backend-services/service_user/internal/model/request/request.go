package request

import (
	"service-user/internal/model"
)

type CreateUserRequest struct {
	ReqID       string         `json:"req_id"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Role        model.RoleType `json:"role"`
}

type UpdateUserRequest struct {
	ReqID       string         `json:"req_id"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	FullName    string         `json:"full_name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Role        model.RoleType `json:"role"`
}

type LoginRequest struct {
	ReqID    string `json:"req_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
