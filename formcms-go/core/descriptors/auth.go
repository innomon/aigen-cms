package descriptors

import (
	"time"
)

const UserTableName = "__users"

const (
	RoleSa    = "sa"
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleGuest = "guest"
)

type User struct {
	Id           int64     `json:"id" mapstructure:"id"`
	Email        string    `json:"email" mapstructure:"email"`
	PasswordHash string    `json:"-" mapstructure:"password_hash"`
	Roles        []string  `json:"roles" mapstructure:"roles"`
	AvatarPath   string    `json:"avatarPath" mapstructure:"avatar_path"`
	CreatedAt    time.Time `json:"createdAt" mapstructure:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" mapstructure:"updated_at"`
}
