package services

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"testing/dating/api/pkg/app"
	"testing/dating/api/pkg/domain"
)

type Services struct {
	db    *gorm.DB
	user  *UserService
	admin *AdminService
}

// createContext - init services and transaction for test, return func for rollback
func createContext(t *testing.T) (*Services, func()) {
	config := app.CreateConfig()

	t.Log("create context")

	lgr := log.Logger
	ctx := context.Background()
	db, err := gorm.Open(postgres.Open(config.Dns), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to open database")
	}

	userService, err := NewUserService(ctx, db, lgr, config.SmsToken, config.Email)
	if err != nil {
		t.Fatal(err)
	}

	adminService, err := NewAdminService(ctx, db, lgr)
	if err != nil {
		t.Fatal(err)
	}

	services := &Services{
		db:    db,
		user:  userService,
		admin: adminService,
	}

	tx := db.Begin()
	tx.SavePoint("sp1")

	services.db = tx
	services.user.SetDb(tx)
	services.admin.SetDb(tx)

	return services, func() {
		if err := tx.RollbackTo("sp1").Error; err != nil {
			panic(err)
		}
	}
}

func equalUsers(t *testing.T, u1, u2 *domain.User) {
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Hash, u2.Hash)
	assert.Equal(t, u1.Phone, u2.Phone)
	assert.Equal(t, u1.FullName, u2.FullName)
}
