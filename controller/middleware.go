package controller

import (
	"errors"
	"fmt"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/content-api/repository"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

var errTokenRequired = errors.New("no access credentials provided")

// GetBearerAuth extracts OAuth access token from request header.
// Authorization: Bearer 19c7d9016b68221cc60f00afca7c498c36c361e3
func GetBearerAuth(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	authForm := req.Form.Get("access_token")

	if authHeader == "" && authForm == "" {
		return "", errTokenRequired
	}

	if authHeader == "" && authForm != "" {
		return authForm, nil
	}

	s := strings.SplitN(authHeader, " ", 2)

	bearerExists := (len(s) == 2) && (strings.ToLower(s[0]) == "bearer")

	log.Printf("Bearer exists: %t", bearerExists)

	if !bearerExists {
		return "", errTokenRequired
	}

	return s[1], nil
}

type AccessGuard struct {
	Env repository.OAuthEnv
}

func (a AccessGuard) CheckToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			_ = view.Render(w, view.NewBadRequest(err.Error()))
			return
		}

		token, err := GetBearerAuth(req)

		if err != nil {
			log.Printf("Token not found: %s", err)

			_ = view.Render(w, view.NewForbidden("Invalid access token"))
			return
		}

		access, err := a.Env.Load(token)

		if err != nil {
			_ = view.Render(w, view.NewDBFailure(err))
			return
		}

		if access.Expired() || !access.Active {
			log.Printf("Token %s is either expired or not active", token)
			_ = view.Render(w, view.NewForbidden("The access token is expired or no longer active"))
			return
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// LogRequest print request headers.
func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		dump, err := httputil.DumpRequest(req, false)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
		log.Printf(string(dump))

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
