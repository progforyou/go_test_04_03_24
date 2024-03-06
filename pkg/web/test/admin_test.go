package test

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/services"
	"testing/dating/api/pkg/web/admin"
)

func setAdminCookie(req *http.Request, sr *Services, t *testing.T) {
	adm := domain.Admin{
		Login:    "admin",
		Password: "admin",
	}
	authorizedAdmin, err := sr.adminService.SignIn(&adm)
	assert.Equal(t, err, nil)
	session := sr.adminSessionController.New(authorizedAdmin)
	req.AddCookie(&http.Cookie{
		Name:     admin.TokenName,
		Value:    session.Token,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   math.MaxInt32,
	})
}

func generateUsers(sr *Services, t *testing.T) []*domain.User {
	user1 := &domain.User{
		Email:    "client@client.ru",
		FullName: "Client client",
		Password: "client",
	}
	err := services.UserServiceUserPrecreated(user1, sr.db)
	assert.Equal(t, err, nil)

	user2 := &domain.User{
		Email:    "client2@client2.ru",
		FullName: "Client2 client2",
		Password: "client2",
	}
	err = services.UserServiceUserPrecreated(user2, sr.db)
	assert.Equal(t, err, nil)
	return []*domain.User{user1, user2}
}

func TestWebAdminSignIn(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	w := httptest.NewRecorder()
	body := &admin.SignInAdmin{
		Login:    "admin",
		Password: "admin",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	assert.Equal(t, err, nil)
	req, _ := http.NewRequest("POST", "/api/v1/admin/auth/signIn", &buf)
	req.Header.Set("Content-Type", "application/json")
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestWebAdminGetUser(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.ru",
		FullName: "Client client",
		Password: "client",
	}
	err := services.UserServiceUserPrecreated(user, sr.db)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/admin/user/%d/", user.ID), nil)
	setAdminCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	foundUser, err := sr.adminService.GetUser(user.ID)
	assert.Equal(t, err, nil)
	equalUsers(t, user, foundUser)
}

func TestWebAdminDeleteUser(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.ru",
		FullName: "Client client",
		Password: "client",
	}
	err := services.UserServiceUserPrecreated(user, sr.db)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/admin/user/%d/", user.ID), nil)
	setAdminCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	_, err = sr.adminService.GetUser(user.ID)
	assert.Equal(t, err, gorm.ErrRecordNotFound)
}

func TestWebAdminGetUsers(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	var res []*domain.User
	list := generateUsers(sr, t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/admin/user/"), nil)
	setAdminCookie(req, sr, t)
	req.URL.Query().Set("page", "0")
	req.URL.Query().Set("per_page", "20")
	sr.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	respBody, err := ioutil.ReadAll(w.Body)
	assert.Equal(t, err, nil)
	err = json.Unmarshal(respBody, &res)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(res), len(list))

	for i, u := range res {
		equalUsers(t, u, list[i])
	}
}

func TestWebAdminUpdateUser(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	user := &domain.User{
		Email:    "client@client.ru",
		FullName: "Client client",
		Password: "client",
	}
	err := services.UserServiceUserPrecreated(user, sr.db)
	assert.Equal(t, err, nil)

	var buf bytes.Buffer
	updates := &admin.UpdateUser{
		FullName: "Client2 client2",
		Phone:    "79343212342",
	}
	err = json.NewEncoder(&buf).Encode(updates)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/admin/user/%d/", user.ID), &buf)
	setAdminCookie(req, sr, t)
	req.Header.Set("Content-Type", "application/json")
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	user.FullName = updates.FullName
	user.Phone = updates.Phone

	foundUser, err := sr.adminService.GetUser(user.ID)
	assert.Equal(t, err, nil)
	equalUsers(t, user, foundUser)
}

func TestWebAdminSignOut(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/admin/auth/signOut", nil)
	setAdminCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}
