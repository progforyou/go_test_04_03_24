package services

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/tools"
)

func TestAdminGetUser(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user1 := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79231234723",
	}
	err := UserServiceUserPrecreated(user1, services.db)
	assert.Equal(t, err, nil)

	user2 := &domain.User{
		Email:    "client2@client2.com",
		FullName: "Client2 client2",
		Password: "password2",
		Phone:    "79231234722",
	}
	err = UserServiceUserPrecreated(user2, services.db)
	assert.Equal(t, err, nil)

	gettingUser, err := services.admin.GetUser(user2.ID)
	assert.Equal(t, err, nil)

	equalUsers(t, user2, gettingUser)

}

func TestAdminGetUsers(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	user1 := &domain.User{
		Email:    "client@client.com",
		FullName: "Client client",
		Password: "password",
		Phone:    "79231234723",
	}
	err := UserServiceUserPrecreated(user1, services.db)
	assert.Equal(t, err, nil)

	user2 := &domain.User{
		Email:    "client2@client2.com",
		FullName: "Client2 client2",
		Password: "password2",
		Phone:    "79231234722",
	}
	err = UserServiceUserPrecreated(user2, services.db)
	assert.Equal(t, err, nil)

	gettingUser, err := services.admin.GetUsers(tools.Page{
		Page:    0,
		PerPage: 20,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(gettingUser), 2)

}

func TestAdminDeleteUser(t *testing.T) {
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

	err = services.admin.DeleteUser(user.ID)
	assert.Equal(t, err, nil)

	_, err = services.admin.GetUser(user.ID)
	assert.NotEqual(t, err, nil)

}

func TestAdminUpdateUser(t *testing.T) {
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

	modifiedUser, err := services.admin.UpdateUser(user)
	assert.Equal(t, err, nil)

	equalUsers(t, modifiedUser, user)

}

func TestAdminSignIn(t *testing.T) {
	services, rollback := createContext(t)
	defer rollback()

	admin := &domain.Admin{
		Login:    "admin",
		Password: "admin",
	}

	err := AdminServiceAdminPrecreated(admin, services.db)
	assert.Equal(t, err, nil)

	signInAdmin, err := services.admin.SignIn(admin)
	assert.Equal(t, err, nil)

	signInAdmin.Login = admin.Login
	signInAdmin.Hash = admin.Hash

}
