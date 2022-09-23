package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/internal/repository"
	"gitlab.com/ftchinese/content-api/pkg/xhttp"
	"go.uber.org/zap"
	"net/http"
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
