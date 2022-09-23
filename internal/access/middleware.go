package access

import (
	"errors"
	"github.com/FTChinese/go-rest/view"
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
			_ = view.Render(w, view.NewBadRequest(err.Error()))
			return
		}

		token, err := GetBearerAuth(req)

		if err != nil {
			log.Printf("Token not found: %s", err)

			_ = view.Render(w, view.NewForbidden("Invalid access token"))
			return
		}

		access, err := a.env.Load(token)

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
