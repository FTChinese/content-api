package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/internal/repository"
	"gitlab.com/ftchinese/content-api/pkg/xhttp"
	"go.uber.org/zap"
	"net/http"
)

type VideoRouter struct {
	repo repository.VideoEnv
}

func NewVideoRouter(db *sqlx.DB, logger *zap.Logger) VideoRouter {
	return VideoRouter{
		repo: repository.NewVideoEnv(db, logger),
	}
}

func (router VideoRouter) Article(w http.ResponseWriter, req *http.Request) {
	id, err := xhttp.GetURLParam(req, "id").ToInt()

	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	video, err := router.repo.Load(id)

	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(video)
}
