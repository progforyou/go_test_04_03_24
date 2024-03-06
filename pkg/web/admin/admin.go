package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
	"math"
	"net/http"
	"strconv"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/tools"
	"testing/dating/api/pkg/web/common"
)

type UserController struct {
	logger            zerolog.Logger
	adminService      domain.AdminService
	sessionController *SessionController
}

func NewAdminUserController(logger zerolog.Logger, sessionController *SessionController, adminService domain.AdminService) *UserController {
	c := &UserController{
		logger:            logger,
		adminService:      adminService,
		sessionController: sessionController,
	}
	return c
}

// GetUser 		        godoc
// @Summary				Get user by id.
// @Param               user_id path int true "User ID"
// @Description			Return user.
// @Produce				application/json
// @Tags				admin
// @Success				200 {object} domain.User{}
// @failure				404 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /admin/user/{user_id} [get]
func (c *UserController) GetUser(g *gin.Context) {
	user := g.MustGet("user").(*domain.User)
	g.JSON(http.StatusOK, user)
}

// GetUsers 		    godoc
// @Summary				Get users.
// @Param               page query int true "Page number" default(0)
// @Param               per_page query int true "Object count in page" default(25)
// @Description			Return users.
// @Produce				application/json
// @Tags				admin
// @Success				200 {object} []domain.User{}
// @failure				404 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /admin/user/{user_id} [get]
func (c *UserController) GetUsers(g *gin.Context) {
	page := tools.GetPage(g)
	users, err := c.adminService.GetUsers(page)
	if err != nil {
		return
	}
	g.JSON(http.StatusOK, users)
}

// UpdateUser 		    godoc
// @Summary				Update user.
// @Param request       body UpdateUser true "query params"
// @Description			Return user.
// @Produce				application/json
// @Tags				admin
// @Success				200 {object} domain.User{}
// @failure				403 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /admin/user/{user_id} [put]
func (c *UserController) UpdateUser(g *gin.Context) {
	var changes UpdateUser
	err := g.ShouldBind(&changes)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	if err = validator.Validate(changes); err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	user := g.MustGet("user").(*domain.User)
	user.FullName = changes.FullName
	user.Phone = changes.Phone
	user, err = c.adminService.UpdateUser(user)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.JSON(http.StatusOK, user)
}

// DeleteUser 		    godoc
// @Summary				Delete user.
// @Description			Return nothing.
// @Produce				application/json
// @Tags				admin
// @Success				200
// @failure				403 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /admin/user/{user_id} [delete]
func (c *UserController) DeleteUser(g *gin.Context) {
	user := g.MustGet("user").(*domain.User)
	err := c.adminService.DeleteUser(user.ID)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.JSON(http.StatusOK, gin.H{})
}

// SignIn 		        godoc
// @Summary				Sign-in user.
// @Param request       body SignInAdmin true "query params"
// @Description			Return session.
// @Produce				application/json
// @Tags				admin
// @Success				200 {object} Session{}
// @failure				404 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /admin/auth/signIn [post]
func (c *UserController) SignIn(g *gin.Context) {
	var admin SignInAdmin
	err := g.ShouldBind(&admin)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.WebError{Err: err.Error()})
		return
	}
	if err = validator.Validate(admin); err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	domainAdmin, err := c.adminService.SignIn(admin.toModel())
	if err != nil {
		common.SendError(g, err)
		return
	}
	session := c.sessionController.New(domainAdmin)
	g.SetCookie(TokenName, session.Token, math.MaxInt32, "/", "", false, true)
	g.JSON(http.StatusOK, Session{
		Token: session.Token,
		Admin: domainAdmin,
	})
}

// SignOut 		        godoc
// @Summary				Remove user session
// @Description			Return Status OK.
// @Produce				application/json
// @Tags				admin
// @Success				200
// @failure				401 {object} common.WebError{}
// @Router              /auth/signOut [get]
func (c *UserController) SignOut(g *gin.Context) {
	token, err := g.Cookie(TokenName)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.WebError{Err: "Unauthorized"})
		return
	}
	log.Debug().Interface("token", token).Msg("SignOut")
	c.sessionController.Remove(token)
	g.SetCookie(TokenName, "", -1, "/", "", false, true)
	g.JSON(http.StatusOK, nil)
}

func (c *UserController) UserMiddleware(g *gin.Context) {
	val, err := strconv.ParseUint(g.Param("user_id"), 10, 64)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
	}
	if val == 0 {
		g.JSON(http.StatusBadRequest, common.WebError{Err: "user_id is required"})
	}
	g.Set("user_id", val)
	u, err := c.adminService.GetUser(val)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.Set("user", u)
	g.Next()
}

type SignInAdmin struct {
	Login    string `form:"login" json:"login" xml:"name"  binding:"required" example:"admin" validate:"min=3,max=120"`
	Password string `form:"password" json:"password" xml:"password" binding:"required" example:"password" validate:"min=3,max=40"`
}

func (u *SignInAdmin) toModel() *domain.Admin {
	return &domain.Admin{
		Login:    u.Login,
		Password: u.Password,
	}
}

type UpdateUser struct {
	FullName string `form:"name" json:"name" binding:"required" example:"John Doe" validate:"min=3,max=120"`
	Phone    string `form:"phone" json:"phone" binding:"required" example:"79123456789" validate:"min=11,max=12"`
}
