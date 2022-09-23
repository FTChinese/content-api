package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"gitlab.com/ftchinese/content-api/internal/repository"
	"gitlab.com/ftchinese/content-api/pkg/xhttp"
	"go.uber.org/zap"
	"net/http"
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
