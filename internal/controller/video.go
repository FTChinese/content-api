package controller

import (
	"net/http"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/repository"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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
