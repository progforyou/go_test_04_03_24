package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/sms"
	"testing/dating/api/pkg/tools"
	"time"
)

var codeTTL = 15 * time.Minute

type UserService struct {
	db     *gorm.DB
	c      *sms.SmsClient
	logger zerolog.Logger
	ctx    context.Context
}

func (u *UserService) Activate(user *domain.User) (*domain.UserActivate, error) {
	var err error
	activate := &domain.UserActivate{}
	activate.UserID = user.ID
	activate.Code = fmt.Sprintf("%d", tools.RandomCode(6))
	err = u.c.SendSms(user.Phone, fmt.Sprintf("Код подтверждения: %s", activate.Code))
	log.Debug().Str("code", activate.Code).Msg("activate")
	if err != nil {
		return nil, err
	}
	err = u.db.Save(activate).Error
	if err != nil {
		return nil, err
	}
	return activate, nil
}

func (u *UserService) ActivateCode(user *domain.User, userActivate *domain.UserActivate) (*domain.User, error) {
	var foundUserActivate domain.UserActivate
	endCode := time.Now().Add(-codeTTL)
	err := u.db.Where("token = ?", userActivate.Token).First(&foundUserActivate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.NewWSError(http.StatusBadRequest, "wrong token")
		}
		return nil, err
	}
	if foundUserActivate.CreatedAt.Before(endCode) {
		return nil, tools.NewWSError(http.StatusBadRequest, "lifetime code done")
	}
	if foundUserActivate.Code != userActivate.Code {
		return nil, tools.NewWSError(http.StatusBadRequest, "wrong code")
	}
	user.Activated = true
	err = u.db.Delete(&foundUserActivate).Error
	if err != nil {
		return nil, err
	}

	return u.Update(user)
}

func (u *UserService) Update(user *domain.User) (*domain.User, error) {
	err := u.db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) SignIn(user *domain.User) (*domain.User, error) {
	var foundUser domain.User
	err := u.db.Where("email = ?", user.Email).First(&foundUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.NewWSError(http.StatusBadRequest, "wrong email")
		}
		return nil, err
	}
	if tools.HashPassword(user.Password) != foundUser.Hash {
		return nil, tools.NewWSError(http.StatusBadRequest, "wrong password")
	}
	return &foundUser, nil
}

func (u *UserService) SignUp(user *domain.User) (*domain.User, error) {
	if u.db.Model(&user).Where("email = ?", user.Email).Updates(&user).RowsAffected == 0 {
		err := u.db.Create(&user).Error
		if err != nil {
			if tools.IsErrUniqueConstraint(err) {
				return nil, tools.NewWSError(http.StatusBadRequest, "email already exist")
			}
		}
		return user, err
	}
	return nil, tools.NewWSError(http.StatusBadRequest, "email already exist")
}

func (u *UserService) SignOut() error {
	return nil
}

func UserServiceUserPrecreated(user *domain.User, db *gorm.DB) error {
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
			log.Debug().Uint64("id", user.ID).Str("email", user.Email).Str("password", user.Password).Msg("User created")
			return nil
		}
		return err
	}
	log.Debug().Uint64("id", user.ID).Str("email", user.Email).Msg("User already exist")
	return nil
}

func (u *UserService) SetDb(tx *gorm.DB) {
	u.db = tx
}

func NewUserService(ctx context.Context, db *gorm.DB, logger zerolog.Logger, smsToken string, email string) (*UserService, error) {
	if err := db.AutoMigrate(&domain.User{}, &domain.UserActivate{}); err != nil {
		return nil, err
	}

	u := UserService{
		ctx:    ctx,
		db:     db,
		logger: logger,
		c:      sms.NewSmsClient(email, smsToken),
	}

	return &u, nil
}
