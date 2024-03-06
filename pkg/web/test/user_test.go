package test

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/dating/api/pkg/domain"
	"testing/dating/api/pkg/web/admin"
	"testing/dating/api/pkg/web/client"
)

func setUserCookie(req *http.Request, sr *Services, t *testing.T) *domain.User {
	session := getUserSession(sr, t)
	req.AddCookie(&http.Cookie{
		Name:     client.TokenName,
		Value:    session.Token,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		MaxAge:   math.MaxInt32,
	})
	return session.User
}

func getUserSession(sr *Services, t *testing.T) *client.Session {
	u := domain.User{
		Email:    "client@client.com",
		Password: "client",
	}
	authorizedUser, err := sr.userService.SignIn(&u)
	assert.Equal(t, err, nil)
	return sr.clientSessionController.New(authorizedUser)
}

func TestWebClientSignIn(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	w := httptest.NewRecorder()
	body := &client.SignInUser{
		EMail:    "client@client.com",
		Password: "client",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	assert.Equal(t, err, nil)
	req, _ := http.NewRequest("POST", "/api/v1/client/auth/signIn", &buf)
	req.Header.Set("Content-Type", "application/json")
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestWebClientSignUp(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	w := httptest.NewRecorder()
	body := &client.SignUpUser{
		SignInUser: client.SignInUser{
			EMail:    "client2@client2.com",
			Password: "client2",
		},
		FullName: "Test test",
		Phone:    "23232322222",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	assert.Equal(t, err, nil)
	req, _ := http.NewRequest("POST", "/api/v1/client/auth/signUp", &buf)
	req.Header.Set("Content-Type", "application/json")
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}

func TestWebClientGet(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()
	var res domain.User

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/client/user/", nil)
	u := setUserCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	respBody, err := ioutil.ReadAll(w.Body)
	assert.Equal(t, err, nil)
	err = json.Unmarshal(respBody, &res)
	assert.Equal(t, err, nil)

	equalUsers(t, u, &res)
}

func TestWebClientUpdate(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	var res domain.User
	var buf bytes.Buffer
	updates := &admin.UpdateUser{
		FullName: "Client2 client2",
		Phone:    "79343212342",
	}
	err := json.NewEncoder(&buf).Encode(updates)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/client/user/", &buf)
	req.Header.Set("Content-Type", "application/json")
	setUserCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	respBody, err := ioutil.ReadAll(w.Body)
	assert.Equal(t, err, nil)
	err = json.Unmarshal(respBody, &res)
	assert.Equal(t, err, nil)
	assert.Equal(t, res.FullName, updates.FullName)
	assert.Equal(t, res.Phone, updates.Phone)
}

func TestWebClientActivateCode(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()
	var activateUser domain.UserActivate

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/client/user/activate", nil)
	_ = setUserCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	respBody, err := ioutil.ReadAll(w.Body)
	assert.Equal(t, err, nil)
	err = json.Unmarshal(respBody, &activateUser)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, activateUser.Token, "")

	err = sr.db.Where("token = ?", activateUser.Token).First(&domain.UserActivate{}).Error
	assert.Equal(t, err, nil)
}

func TestWebClientActivate(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	code := "111111"
	session := getUserSession(sr, t)
	user := session.User
	codeInBase := &domain.UserActivate{
		UserID: user.ID,
		Code:   code,
	}
	sr.db.Save(codeInBase)

	var buf bytes.Buffer
	activateUser := &client.ActivateUser{
		Token: codeInBase.Token,
		Code:  code,
	}
	err := json.NewEncoder(&buf).Encode(activateUser)
	assert.Equal(t, err, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/client/user/activate", &buf)
	setUserCookie(req, sr, t)
	req.Header.Set("Content-Type", "application/json")
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	foundUser, err := sr.adminService.GetUser(user.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, foundUser.Activated, true)
}

func TestWebClientSignOut(t *testing.T) {
	sr, rollback := createContext(t)
	defer rollback()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/client/auth/signOut", nil)
	setUserCookie(req, sr, t)
	sr.router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
}
