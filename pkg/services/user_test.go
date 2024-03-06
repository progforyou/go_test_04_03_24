package services

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"testing/dating/api/pkg/domain"
	"time"
)

func TestClientActivateCodeUser(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79281724695",
	}
	err := UserServiceUserPrecreated(user, services.db)
	assert.Equal(t, err, nil)

	userActivate := &domain.UserActivate{
		UserID: user.ID,
		Code:   "111111",
	}
	err = services.db.Save(userActivate).Error
	assert.Equal(t, err, nil)

	time.Sleep(1 * time.Second)

	activatedUser, err := services.user.ActivateCode(user, userActivate)
	assert.Equal(t, err, nil)
	assert.Equal(t, activatedUser.Activated, true)

}

func TestClientActivateUser(t *testing.T) {
	t.Skip()

	services, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79281724695",
	}
	err := UserServiceUserPrecreated(user, services.db)
	assert.Equal(t, err, nil)

	activateUser, err := services.user.Activate(user)
	assert.Equal(t, err, nil)
	assert.Equal(t, activateUser.UserID, user.ID)
}

func TestClientUpdateUser(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79231234723",
	}
	err := UserServiceUserPrecreated(user, services.db)
	assert.Equal(t, err, nil)

	user.FullName = "Client updated"
	user.Email = "client_new@client_new.com"

	modifiedUser, err := services.user.Update(user)
	assert.Equal(t, err, nil)

	equalUsers(t, modifiedUser, user)

}

func TestClientSignIn(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.com",
		Password: "client",
	}

	err := UserServiceUserPrecreated(user, services.db)
	assert.Equal(t, err, nil)

	signInAdmin, err := services.user.SignIn(user)
	assert.Equal(t, err, nil)

	signInAdmin.Email = user.Email
	signInAdmin.Hash = user.Hash

}

func TestClientSignUp(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79231234723",
	}

	createdUser, err := services.user.SignUp(user)
	assert.Equal(t, err, nil)

	assert.Equal(t, createdUser.Email, user.Email)
	assert.Equal(t, createdUser.FullName, user.FullName)
	assert.Equal(t, createdUser.Hash, user.Hash)
	assert.Equal(t, createdUser.Phone, user.Phone)

}
