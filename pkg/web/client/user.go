package client

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
	"math"
	"net/http"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/web/common"
)

type UserController struct {
	logger            zerolog.Logger
	userService       domain.UserService
	sessionController *SessionController
}

func NewClientUserController(logger zerolog.Logger, sessionController *SessionController, userService domain.UserService) *UserController {
	c := &UserController{
		logger:            logger,
		userService:       userService,
		sessionController: sessionController,
	}
	return c
}

// Get 		            godoc
// @Summary				Get user.
// @Description			Return user.
// @Produce				application/json
// @Tags				user
// @Success				200 {object} domain.User{}
// @failure				404 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /user/ [get]
func (c *UserController) Get(g *gin.Context) {
	user := g.MustGet("user").(*domain.User)
	g.JSON(http.StatusOK, user)
}

// Update 		        godoc
// @Summary				Update user.
// @Param request       body UpdateUser true "query params"
// @Description			Return user.
// @Produce				application/json
// @Tags				user
// @Success				200 {object} domain.User{}
// @failure				403 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /user/ [put]
func (c *UserController) Update(g *gin.Context) {
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

	user, err = c.userService.Update(user)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.JSON(http.StatusOK, user)
}

// Activate 		    godoc
// @Summary				Activate user.
// @Description			Return token.
// @Produce				application/json
// @Tags				user
// @Success				200 {object} domain.UserActivate{}
// @failure				403 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /user/activate [get]
func (c *UserController) Activate(g *gin.Context) {
	user := g.MustGet("user").(*domain.User)
	activate, err := c.userService.Activate(user)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.JSON(http.StatusOK, activate)
}

// ActivateCode 		godoc
// @Summary				ActivateCode user.
// @Description			Return token.
// @Produce				application/json
// @Tags				user
// @Success				200 {object} domain.User{}
// @failure				403 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /user/activate [post]
func (c *UserController) ActivateCode(g *gin.Context) {
	user := g.MustGet("user").(*domain.User)
	var activateUserReq ActivateUser
	err := g.ShouldBind(&activateUserReq)
	if err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	if err = validator.Validate(activateUserReq); err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	activateUser := activateUserReq.toModel()
	activateUser.UserID = user.ID
	modifiedUser, err := c.userService.ActivateCode(user, activateUser)
	if err != nil {
		common.SendError(g, err)
		return
	}
	g.JSON(http.StatusOK, modifiedUser)
}

// SignIn 		        godoc
// @Summary				Sign-in user.
// @Param request       body SignInUser true "query params"
// @Description			Return session.
// @Produce				application/json
// @Tags				auth
// @Success				200 {object} Session{}
// @failure				404 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /auth/signIn [post]
func (c *UserController) SignIn(g *gin.Context) {
	var user SignInUser
	err := g.ShouldBind(&user)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.WebError{Err: err.Error()})
		return
	}
	if err = validator.Validate(user); err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}
	domainUser, err := c.userService.SignIn(user.toModel())
	if err != nil {
		common.SendError(g, err)
		return
	}
	session := c.sessionController.New(domainUser)
	g.SetCookie(TokenName, session.Token, math.MaxInt32, "/", "", false, true)
	g.JSON(http.StatusOK, Session{
		Token: session.Token,
		User:  domainUser,
	})
}

// SignUp 		        godoc
// @Summary				Sign-up new user.
// @Param request       body SignUpUser{} true "query params"
// @Description			Return session.
// @Produce				application/json
// @Tags				auth
// @Success				200 {object} Session{}
// @failure				402 {object} common.WebError{}
// @failure				401 {object} common.WebError{}
// @Router              /auth/signUp [post]
func (c *UserController) SignUp(g *gin.Context) {
	var user SignUpUser
	err := g.ShouldBind(&user)
	if err != nil {
		g.JSON(http.StatusUnauthorized, common.WebError{Err: err.Error()})
		return
	}
	if err = validator.Validate(user); err != nil {
		g.JSON(http.StatusBadRequest, common.WebError{Err: err.Error()})
		return
	}

	domainUser, err := c.userService.SignUp(user.toModel())
	if err != nil {
		common.SendError(g, err)
		return
	}

	log.Debug().Interface("domainUser", domainUser).Msg("SignUp")
	session := c.sessionController.New(domainUser)
	log.Debug().Interface("session", session).Msg("SignUp")
	g.SetCookie(TokenName, session.Token, math.MaxInt32, "/", "", false, true)
	g.JSON(http.StatusOK, Session{
		Token: session.Token,
		User:  domainUser,
	})
}

// SignOut 		        godoc
// @Summary				Remove user session
// @Description			Return Status OK.
// @Produce				application/json
// @Tags				auth
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

type SignInUser struct {
	EMail    string `form:"email" json:"email" xml:"name"  binding:"required" example:"example@example.com" validate:"min=3,max=120,regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password string `form:"password" json:"password" xml:"password" binding:"required" example:"password" validate:"min=3,max=40"`
}

func (u *SignInUser) toModel() *domain.User {
	return &domain.User{
		Email:    u.EMail,
		Password: u.Password,
	}
}

type SignUpUser struct {
	SignInUser
	FullName string `form:"name" json:"name" binding:"required" example:"John Doe" validate:"min=3,max=120"`
	Phone    string `form:"phone" json:"phone" binding:"required" example:"79123456789" validate:"min=11,max=12"`
}

func (u *SignUpUser) toModel() *domain.User {
	return &domain.User{
		Email:    u.EMail,
		Password: u.Password,
		FullName: u.FullName,
		Phone:    u.Phone,
	}
}

type UpdateUser struct {
	FullName string `form:"name" json:"name" binding:"required" example:"John Doe" validate:"min=3,max=120"`
	Phone    string `form:"phone" json:"phone" binding:"required" example:"79123456789" validate:"min=11,max=12"`
}

type ActivateUser struct {
	Token string `form:"token" json:"token" binding:"required" example:"-"`
	Code  string `form:"code" json:"code" binding:"required" example:"111111" validate:"min=6,max=6"`
}

func (u *ActivateUser) toModel() *domain.UserActivate {
	return &domain.UserActivate{
		Token: u.Token,
		Code:  u.Code,
	}
}
