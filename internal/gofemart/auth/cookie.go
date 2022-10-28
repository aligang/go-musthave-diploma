package auth

import (
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
	"sync"
	"time"
)

type cookieInfo struct {
	login     string
	validTill time.Time
}

type cookiesMeta map[string]cookieInfo

type Auth struct {
	cookiesMeta cookiesMeta
	Lock        sync.Mutex
}

func New() *Auth {
	auth := &Auth{cookiesMeta: cookiesMeta{}}
	go auth.junkCleaner()
	return auth
}

func (a *Auth) CreateAuthCookie(account *account.CustomerAccount) *http.Cookie {
	logging.Debug("Creating New cookie for %s", account.Login)
	token := account.Login + account.Password
	a.cookiesMeta[token] = cookieInfo{
		login:     account.Login,
		validTill: time.Now().Add(time.Second * time.Duration(CookieLifetime)),
	}
	return &http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: CookieLifetime,
	}
}

func (a *Auth) CreateUserInfoCookie(token string) *http.Cookie {
	metadata := a.cookiesMeta[token]
	return &http.Cookie{
		Name:  "username",
		Value: metadata.login,
	}
}

func (a *Auth) CheckAuthCookie(cookie *http.Cookie) error {
	meta, ok := a.cookiesMeta[cookie.Value]
	if !ok || time.Now().After(meta.validTill) {
		return fmt.Errorf("cookie is invalid or expired")
	}
	return nil
}
