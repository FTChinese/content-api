package controller

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type GalleryRouter struct {
	env repository.Env
}

func NewGalleryStory(db *sqlx.DB) GalleryRouter {
	return GalleryRouter{env: repository.NewEnv(db)}
}

func (router GalleryRouter) Article(w http.ResponseWriter, req *http.Request) {

}
