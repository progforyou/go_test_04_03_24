package app

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"testing/dating/api/pkg/domain"
)

func Initial(ctx context.Context, db *gorm.DB, logger zerolog.Logger, userService domain.UserService, adminService domain.AdminService) error {
	var admins []*domain.Admin
	var users []*domain.User
	err := db.Find(&admins).Error
	if err != nil {
		return err
	}
	if len(admins) == 0 {
		//Create admin user
		admin := domain.Admin{
			Password: "admin",
			Login:    "admin",
		}
		if err = db.Create(&admin).Error; err != nil {
			return err
		}
		log.Info().Str("login", admin.Login).Str("password", admin.Password).Msg("Admin user created")
	}
	err = db.Find(&users).Error
	if err != nil {
		return err
	}
	if len(users) == 0 {
		//Create admin user
		user := domain.User{
			Password: "client",
			Email:    "client@client.com",
		}
		if err = db.Create(&user).Error; err != nil {
			return err
		}
		log.Info().Str("Email", user.Email).Str("password", user.Password).Msg("Client user created")
	}
	return nil
}
