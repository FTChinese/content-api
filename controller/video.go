package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type VideoRouter struct {
	env repository.Env
}

func NewVideoRouter(db *sqlx.DB) VideoRouter {
	return VideoRouter{env: repository.NewEnv(db)}
}

func (router VideoRouter) Article(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	video, err := router.env.RetrieveVideo(id)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(video))
}
