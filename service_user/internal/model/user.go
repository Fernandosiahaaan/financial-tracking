package model

import "time"

var (
	UserSessionTime = 1 * time.Hour
)

type User struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type LoginData struct {
	Id    string `json:"user_id"`
	Token string `json:"token"`
}

type LoginCacheData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
