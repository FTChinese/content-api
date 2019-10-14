package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type ChannelRouter struct {
	env repository.Env
}

func NewChannelRouter(db *sqlx.DB) ChannelRouter {
	return ChannelRouter{env: repository.NewEnv(db)}
}

func (router ChannelRouter) TodayFrontPage(w http.ResponseWriter, req *http.Request) {
	frontPage, err := router.env.TodayFrontPage()

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(frontPage))
}

func (router ChannelRouter) ArchivedFrontPage(w http.ResponseWriter, req *http.Request) {
	date, err := GetURLParam(req, "date").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	frontPage, err := router.env.ArchivedFrontPage(date)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(frontPage))
}
