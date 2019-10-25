package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type GalleryRouter struct {
	env repository.ContentEnv
}

func NewGalleryStory(db *sqlx.DB) GalleryRouter {
	return GalleryRouter{env: repository.NewContentEnv(db)}
}

func (router GalleryRouter) Article(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	g, err := router.env.RetrieveGallery(id)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(g))
}
