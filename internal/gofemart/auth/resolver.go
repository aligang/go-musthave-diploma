package auth

import (
	"net/http"
)

func ResolveUsername(r *http.Request) (string, error) {
	biscuit, err := r.Cookie("username")
	if err != nil {
		return "", err
	}
	return biscuit.Value, nil
}
