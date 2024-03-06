package domain

import (
	"gorm.io/gorm"
	"testing/dating/api/pkg/tools"
	"time"
)

type AdminService interface {
	SignIn(admin *Admin) (*Admin, error)
	SignOut() error
	GetUser(userId uint64) (*User, error)
	GetUsers(page tools.Page) ([]*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(userId uint64) error
}

type Admin struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Login     string    `gorm:"not null;" json:"login"`
	Hash      string    `gorm:"not null;size:64" json:"-"`
	Password  string    `gorm:"-" json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *Admin) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		u.Hash = tools.HashPassword(u.Password)
	}
	return nil
}
