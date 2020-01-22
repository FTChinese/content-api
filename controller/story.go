package controller

import (
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type StoryRouter struct {
	repo repository.StoryEnv
}

func NewStoryRouter(db *sqlx.DB) StoryRouter {
	return StoryRouter{
		repo: repository.NewStoryEnv(db),
	}
}

func (router StoryRouter) Raw(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.JSON(w, view.NewBadRequest(err.Error()))
		return
	}

	story, err := router.repo.LoadRawStory(id)
	if err != nil {
		_ = view.JSON(w, view.NewDBFailure(err))
		return
	}

	_ = view.JSON(w, view.NewResponse().SetBody(story))
}

func (router StoryRouter) CN(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.JSON(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.repo.LoadRawStory(id)
	if err != nil {
		_ = view.JSON(w, view.NewDBFailure(err))
		return
	}

	_ = view.JSON(w, view.NewResponse().SetBody(rawStory.BuildCN()))
}

func (router StoryRouter) EN(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.JSON(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.repo.LoadRawStory(id)
	if err != nil {
		_ = view.JSON(w, view.NewDBFailure(err))
		return
	}

	s, err := rawStory.BuildEN()
	if err != nil {
		_ = view.JSON(w, view.NewNotFound())
		return
	}

	_ = view.JSON(w, view.NewResponse().SetBody(s))
}

// Bilingual output the article in bilingual format.
func (router StoryRouter) Bilingual(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.JSON(w, view.NewBadRequest(err.Error()))
		return
	}

	rawStory, err := router.repo.RetrieveRawStory(id)
	if err != nil {
		_ = view.JSON(w, view.NewDBFailure(err))
		return
	}

	_ = view.JSON(w, view.NewResponse().SetBody(rawStory.BuildBilingual()))
}
