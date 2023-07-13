package controller

import (
	"net/http"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/pkg"
	"github.com/FTchinese/content-api/internal/repository"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type StoryRouter struct {
	repo   repository.StoryEnv
	logger *zap.Logger
}

func NewStoryRouter(db *sqlx.DB, l *zap.Logger) StoryRouter {
	return StoryRouter{
		repo:   repository.NewStoryEnv(db, l),
		logger: l,
	}
}

func (router StoryRouter) Story(w http.ResponseWriter, req *http.Request) {
	id, err := xhttp.GetURLParam(req, "id").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	story, err := router.repo.LoadRawStory(id)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(pkg.NewStory(story))
}

func (router StoryRouter) StoryNoCache(w http.ResponseWriter, req *http.Request) {
	id, err := xhttp.GetURLParam(req, "id").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	story, err := router.repo.RetrieveStory(id)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(story)
}
