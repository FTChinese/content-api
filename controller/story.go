package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type StoryRouter struct {
	env repository.Env
}

func NewStoryRouter(db *sqlx.DB) StoryRouter {
	return StoryRouter{
		env: repository.NewEnv(db),
	}
}

func (router StoryRouter) Raw(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	story, err := router.env.RetrieveRawStory(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(story))
}
