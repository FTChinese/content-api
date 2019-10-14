package controller

import (
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/content-api/models"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type AccessGuard struct {
	Env repository.OAuthEnv
}

func (a AccessGuard) CheckToken(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			_ = view.Render(w, view.NewBadRequest(err.Error()))
			return
		}

		token, err := models.GetBearerAuth(req)

		if err != nil {
			_ = view.Render(w, view.NewForbidden("Invalid access token"))
			return
		}

		access, err := a.Env.Load(token)

		if err != nil {
			_ = view.Render(w, view.NewDBFailure(err))
			return
		}

		if access.Expired() || !access.Active {
			_ = view.Render(w, view.NewForbidden("The access token is expired or no longer active"))
			return
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
