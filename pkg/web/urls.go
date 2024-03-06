package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/web/admin"
	"testing/dating/api/pkg/web/client"
	"testing/dating/api/pkg/web/docs"
)

//go:generate swag init -g web/urls.go -d ../ -o ./docs --parseDependency

// NewController create new controller
// @title           Dating API
// @version         1.0
// @description     Dating api backend server
// @host            localhost:8084
// @BasePath        /api/v1
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
func Urls(
	ctx context.Context,
	logger zerolog.Logger,
	r *gin.Engine,
	host string,
	userService domain.UserService,
	adminService domain.AdminService,
	clientSessionController *client.SessionController,
	adminSessionController *admin.SessionController,
) (*gin.Engine, error) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = host
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := r.Group("/api/v1")
	{
		cl := v1.Group("/client")
		{
			us := client.NewClientUserController(logger, clientSessionController, userService)
			ag := cl.Group("/auth")
			{
				ag.POST("/signIn", us.SignIn)
				ag.POST("/signUp", us.SignUp)
				ag.GET("/signOut", us.SignOut)
			}
			ug := cl.Group("/user")
			{
				ug.Use(clientSessionController.AuthMiddleware())
				ug.GET("/", us.Get)
				ug.PUT("/", us.Update)
				ug.GET("/activate", us.Activate)
				ug.POST("/activate", us.ActivateCode)
			}
		}

		a := v1.Group("/admin")
		{
			us := admin.NewAdminUserController(logger, adminSessionController, adminService)
			ag := a.Group("/auth")
			{
				ag.POST("/signIn", us.SignIn)
				ag.GET("/signOut", us.SignOut)
			}
			ug := a.Group("/user")
			{
				ug.Use(adminSessionController.AuthMiddleware())
				ug.GET("/", us.GetUsers)

				ugo := ug.Group("/:user_id")
				{
					ugo.Use(us.UserMiddleware)
					ugo.GET("/", us.GetUser)
					ugo.PUT("/", us.UpdateUser)
					ugo.DELETE("/", us.DeleteUser)
				}

			}
		}
	}
	return r, nil
}
