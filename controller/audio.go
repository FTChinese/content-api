package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/models"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type AudioRouter struct {
	env repository.InteractiveEnv
}

func NewAudioRouter(db *sqlx.DB) AudioRouter {
	return AudioRouter{
		env: repository.NewInteractiveEnv(db),
	}
}

func (router AudioRouter) ChannelPage(w http.ResponseWriter, req *http.Request) {
	name, err := GetURLParam(req, "name").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	config, ok := repository.GetAudioChannelConfig(name)

	if !ok {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	p := gorest.GetPagination(req)

	if config.KeyWords.IsZero() {
		_ = view.Render(w, view.NewNotFound())
	}

	teasers, err := router.env.RetrieveChannelTeasers(config.KeyWords.String, p)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(models.ChannelPage{
		ChannelSetting: config,
		Data:           teasers,
	}))
}

func (router AudioRouter) Content(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	article, err := router.env.RetrieveAudioArticle(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(article))
}

func (router AudioRouter) SpeedReadingChannel(w http.ResponseWriter, req *http.Request) {

	p := gorest.GetPagination(req)

	teasers, err := router.env.RetrieveSpeedReadingTeasers(p)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(models.ChannelPage{
		ChannelSetting: repository.SpeedReadingChannel,
		Data:           teasers,
	}))
}

func (router AudioRouter) SpeedReadingArticle(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	article, err := router.env.RetrieveSpeedReading(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(article))
}
