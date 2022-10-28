package auth

import (
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/logging"
	"net/http"
)

func (a *Auth) CheckAuthInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Debug("Authentication: Checking Authentication info")
		biscuit, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "problem during checking auth info: authentication info was not provided", http.StatusUnauthorized)
			logging.Debug("problem during checking auth info: authentication info was not provided")
			return
		}
		err = a.CheckAuthCookie(biscuit)
		if err != nil {
			msg := fmt.Sprintf("problem during checking auth info: %s", err.Error())
			http.Error(w, msg, http.StatusUnauthorized)
			logging.Debug(msg)
			return
		}
		r.AddCookie(a.CreateUserInfoCookie(biscuit.Value))
		logging.Debug("Authentication: Succeeded")
		next.ServeHTTP(w, r)
	})
}
