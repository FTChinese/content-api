package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type VideoRouter struct {
	repo repository.VideoEnv
}

func NewVideoRouter(db *sqlx.DB) VideoRouter {
	return VideoRouter{
		repo: repository.NewVideoEnv(db),
	}
}

func (router VideoRouter) Article(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	video, err := router.repo.Load(id)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(video))
}
