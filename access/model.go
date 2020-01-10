package access

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"log"
	"net/http"
	"strings"
	"time"
)

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

type OAuthAccess struct {
	Token     string      `db:"access_token"`
	Active    bool        `db:"is_active"`
	ExpiresIn null.Int    `db:"expires_in"` // seconds
	CreatedAt chrono.Time `db:"created_utc"`
}

func (o OAuthAccess) Expired() bool {

	if o.ExpiresIn.IsZero() {
		return false
	}

	expireAt := o.CreatedAt.Add(time.Second * time.Duration(o.ExpiresIn.Int64))

	if expireAt.Before(time.Now()) {
		return true
	}

	return false
}