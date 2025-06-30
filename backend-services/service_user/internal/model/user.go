package model

import "time"

var (
	UserSessionTime = 1 * time.Hour
)

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        RoleType  `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoginData struct {
	Id    string `json:"user_id"`
	Token string `json:"token"`
}

type LoginCacheData struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Role     RoleType `json:"role"`
}

type RoleType string

const (
	RoleUser       RoleType = "USER"
	RoleAdmin      RoleType = "ADMIN"
	RoleSuperAdmin RoleType = "SUPERADMIN"
)

func (r RoleType) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin, RoleSuperAdmin:
		return true
	}

	return false
}
