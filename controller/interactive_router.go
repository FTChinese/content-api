package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/models"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type InteractiveRouter struct {
	env repository.InteractiveEnv
}

func NewAudioRouter(db *sqlx.DB) InteractiveRouter {
	return InteractiveRouter{
		env: repository.NewInteractiveEnv(db),
	}
}

func (router InteractiveRouter) ChannelPage(w http.ResponseWriter, req *http.Request) {
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

	teasers, err := router.env.RetrieveTeasers(config.KeyWords.String, p)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	var data []models.InteractiveTeaser
	if config.ContentKind == models.ContentKindSpeedReading {

		for _, v := range teasers {
			data = append(data, v.SpeedReadingTeaser())
		}
	} else {
		for _, v := range teasers {
			data = append(data, v.AudioTeaser())
		}
	}

	_ = view.Render(w, view.NewResponse().SetBody(models.InteractivePage{
		ChannelSetting: config,
		Data:           data,
	}))
}

func (router InteractiveRouter) Content(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	content, err := router.env.RetrieveRawContent(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(
		w,
		view.NewResponse().SetBody(
			content.Build(),
		),
	)
}