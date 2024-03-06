package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"testing/dating/api/pkg/app"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/services"
	"testing/dating/api/pkg/web"
	"testing/dating/api/pkg/web/admin"
	"testing/dating/api/pkg/web/client"
)

type Services struct {
	db                      *gorm.DB
	router                  *gin.Engine
	userService             *services.UserService
	adminService            *services.AdminService
	client                  *http.Client
	adminSessionController  *admin.SessionController
	clientSessionController *client.SessionController
	g                       *gin.Context
}

// createContext - init services and transaction for test, return func for rollback
func createContext(t *testing.T) (*Services, func()) {
	config := app.CreateConfig()
	ctx := context.Background()

	t.Log("create context")

	lgr := log.Logger
	db, err := gorm.Open(postgres.Open(config.Dns), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to open database")
	}

	userService, err := services.NewUserService(ctx, db, lgr, config.SmsToken, config.Email)
	if err != nil {
		t.Fatal(err)
	}

	adminService, err := services.NewAdminService(ctx, db, lgr)
	if err != nil {
		t.Fatal(err)
	}

	s := &Services{
		db:           db,
		userService:  userService,
		adminService: adminService,
	}

	s.g = &gin.Context{
		Request: &http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "localhost",
			},
		},
	}

	s.adminSessionController = admin.NewSessionController(ctx)
	s.clientSessionController = client.NewSessionController(ctx)

	tx := db.Begin()
	tx.SavePoint("sp1")

	s.db = tx
	s.userService.SetDb(tx)
	s.adminService.SetDb(tx)

	err = app.Initial(ctx, s.db, lgr, s.userService, adminService)
	assert.Equal(t, err, nil)

	r := gin.Default()
	r, err = web.Urls(ctx, log.Logger, r, fmt.Sprintf("%s:%s", config.HostName, config.Port), s.userService, s.adminService, s.clientSessionController, s.adminSessionController)
	assert.Equal(t, err, nil)
	s.router = r

	jar, err := cookiejar.New(nil)
	assert.Equal(t, err, nil)

	s.client = &http.Client{
		Jar: jar,
	}

	return s, func() {
		if err := tx.RollbackTo("sp1").Error; err != nil {
			panic(err)
		}
	}
}

func equalUsers(t *testing.T, u1, u2 *domain.User) {
	assert.Equal(t, u1.Email, u2.Email)
	assert.Equal(t, u1.Phone, u2.Phone)
	assert.Equal(t, u1.FullName, u2.FullName)
}
