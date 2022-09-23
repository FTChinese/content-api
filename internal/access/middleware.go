package access

import (
	"errors"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

var errTokenRequired = errors.New("no access credentials provided")

type Guard struct {
	env Repo
}

func NewGuard(db *sqlx.DB) Guard {
	return Guard{
		env: NewRepo(db),
	}
}

func (a Guard) CheckToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			_ = render.New(w).BadRequest(err.Error())
			return
		}

		token, err := GetBearerAuth(req)

		if err != nil {
			log.Printf("Token not found: %s", err)

			_ = render.New(w).Forbidden("Invalid access token")
			return
		}

		access, err := a.env.Load(token)

		if err != nil {
			_ = render.New(w).DBError(err)
			return
		}

		if access.Expired() || !access.Active {
			log.Printf("Token %s is either expired or not active", token)
			_ = render.New(w).Forbidden("The access token is expired or no longer active")
			return
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
