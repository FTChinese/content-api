package controller

import (
	"net/http"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/repository"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type GalleryRouter struct {
	repo repository.GalleryEnv
}

func NewGalleryStory(db *sqlx.DB, logger *zap.Logger) GalleryRouter {
	return GalleryRouter{
		repo: repository.NewGalleryEnv(db, logger),
	}
}

func (router GalleryRouter) Article(w http.ResponseWriter, req *http.Request) {
	id, err := xhttp.GetURLParam(req, "id").ToString()

	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	g, err := router.repo.LoadGallery(id)

	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(g)
}
