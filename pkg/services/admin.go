package services

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/tools"
)

type AdminService struct {
	db     *gorm.DB
	logger zerolog.Logger
	ctx    context.Context
}

func (u *AdminService) UpdateUser(user *domain.User) (*domain.User, error) {
	err := u.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *AdminService) DeleteUser(userId uint64) error {
	err := u.db.Where("id = ?", userId).Delete(&domain.User{ID: userId}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *AdminService) GetUser(userId uint64) (*domain.User, error) {
	var user *domain.User
	err := u.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *AdminService) GetUsers(page tools.Page) ([]*domain.User, error) {
	var list []*domain.User
	err := u.db.Offset(page.Page * page.PerPage).Limit(page.PerPage).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *AdminService) SignIn(admin *domain.Admin) (*domain.Admin, error) {
	var foundAdmin domain.Admin
	err := u.db.Where("login = ?", admin.Login).First(&foundAdmin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.NewWSError(http.StatusBadRequest, "wrong login")
		}
		return nil, err
	}
	if tools.HashPassword(admin.Password) != foundAdmin.Hash {
		return nil, tools.NewWSError(http.StatusBadRequest, "wrong password")
	}
	return &foundAdmin, nil
}

func AdminServiceAdminPrecreated(admin *domain.Admin, db *gorm.DB) error {
	if err := db.Where("login = ?", admin.Login).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&admin).Error; err != nil {
				return err
			}
			log.Debug().Uint64("id", admin.ID).Str("login", admin.Login).Str("password", admin.Password).Msg("Admin created")
			return nil
		}
		return err
	}
	log.Debug().Uint64("id", admin.ID).Str("login", admin.Login).Msg("Admin already exist")
	return nil
}

func (u *AdminService) SignOut() error {
	return nil
}

func (u *AdminService) SetDb(tx *gorm.DB) {
	u.db = tx
}

func NewAdminService(ctx context.Context, db *gorm.DB, logger zerolog.Logger) (*AdminService, error) {
	if err := db.AutoMigrate(&domain.Admin{}); err != nil {
		return nil, err
	}
	u := AdminService{
		ctx:    ctx,
		db:     db,
		logger: logger,
	}
	return &u, nil
}
