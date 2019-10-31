package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/models"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type StoryRouter struct {
	env repository.ContentEnv
}

func NewStoryRouter(db *sqlx.DB) StoryRouter {
	return StoryRouter{
		env: repository.NewContentEnv(db),
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

func (router StoryRouter) CN(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.env.RetrieveRawStory(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(models.NewStoryCN(&rawStory)))
}

func (router StoryRouter) EN(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.env.RetrieveRawStory(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	s, err := models.NewStoryEN(&rawStory)
	if err != nil {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(s))
}

func (router StoryRouter) Bilingual(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.env.RetrieveRawStory(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	s, err := models.NewBilingualStory(&rawStory)
	if err != nil {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(s))
}
