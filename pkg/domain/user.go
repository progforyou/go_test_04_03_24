package domain

import (
	"gorm.io/gorm"
	"testing/dating/api/pkg/tools"
	"time"
)

type UserService interface {
	SignUp(user *User) (*User, error)
	SignIn(user *User) (*User, error)
	SignOut() error
	Update(user *User) (*User, error)
	Activate(user *User) (*UserActivate, error)
	ActivateCode(user *User, userActivate *UserActivate) (*User, error)
}

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"size:120;index:idx_user_mail,unique,size:120;" json:"email"`
	FullName  string    `gorm:"not null;size:120" json:"name"`
	Hash      string    `gorm:"not null;size:64" json:"-"`
	Password  string    `gorm:"-" json:"password,omitempty"`
	Phone     string    `gorm:"default:''" json:"phone"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

type UserActivate struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `gorm:"size:64;not null;index:ids_activate,size:64" json:"token"`
	UserID    uint64    `gorm:"index:ids_activate" json:"-"`
	Code      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (ua *UserActivate) BeforeSave(tx *gorm.DB) error {
	var err error
	ua.Token, err = tools.GenerateToken()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		u.Hash = tools.HashPassword(u.Password)
	}
	u.Email = tools.NormalizeEmail(u.Email)
	return nil
}
