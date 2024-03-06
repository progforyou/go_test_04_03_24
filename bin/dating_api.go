package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing/dating/api/pkg/app"
	"testing/dating/api/pkg/scan"
	"testing/dating/api/pkg/services"
	"testing/dating/api/pkg/web"
	"testing/dating/api/pkg/web/admin"
	"testing/dating/api/pkg/web/client"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.DebugLevel)
	config := app.CreateConfig()

	db, err := gorm.Open(postgres.Open(config.Dns), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("fail to open database")
	}

	ctx := context.Background()
	userService, err := services.NewUserService(ctx, db, log.Logger, config.SmsToken, config.Email)
	if err != nil {
		log.Fatal().Err(err).Msg("fail create UserService")
	}
	adminService, err := services.NewAdminService(ctx, db, log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("fail create AdminService")
	}
	err = app.Initial(ctx, db, log.Logger, userService, adminService)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to initial app")
	}

	clientSessionController := client.NewSessionController(ctx)
	adminSessionController := admin.NewSessionController(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("fail create AdminService")
	}

	//scanner
	scanHandler, err := scan.NewScanHandler(ctx, db, log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to start scan handler")
	}
	scanHandler.Start()
	defer scanHandler.Stop()

	r := gin.New()
	r.Use(gin.Recovery())

	swaggerHost := fmt.Sprintf("%s:%s", config.HostName, config.Port)
	r, err = web.Urls(ctx, log.Logger, r, swaggerHost, userService, adminService, clientSessionController, adminSessionController)
	if err != nil {
		log.Fatal().Err(err).Msg("fail create web")
	}
	bind := config.Host + ":" + config.Port
	log.Debug().Msgf("Start application on port %s", config.Port)
	err = r.Run(bind)
	if err != nil {
		log.Fatal().Err(err).Msg("Server at " + bind)
	}
}
