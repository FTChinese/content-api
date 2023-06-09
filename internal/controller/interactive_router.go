package controller

import (
	"net/http"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/pkg"
	"github.com/FTchinese/content-api/internal/repository"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type InteractiveRouter struct {
	env repository.InteractiveEnv
}

func NewInteractiveRouter(db *sqlx.DB, logger *zap.Logger) InteractiveRouter {
	return InteractiveRouter{
		env: repository.NewInteractiveEnv(db, logger),
	}
}

func (router InteractiveRouter) ChannelPage(w http.ResponseWriter, req *http.Request) {
	name, err := xhttp.GetURLParam(req, "name").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	config, ok := repository.GetAudioChannelConfig(name)

	if !ok {
		_ = render.New(w).NotFound("")
		return
	}

	p := gorest.GetPagination(req)

	if config.KeyWords.IsZero() {
		_ = render.New(w).NotFound("")
		return
	}

	teasers, err := router.env.RetrieveTeasers(config.KeyWords.String, p)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	var data []pkg.Teaser
	for _, v := range teasers {
		data = append(data, v.Teaser())
	}

	_ = render.New(w).OK(pkg.ChannelPage{
		ChannelSetting: config,
		Data:           data,
	})
}

func (router InteractiveRouter) Content(w http.ResponseWriter, req *http.Request) {
	id, err := xhttp.GetURLParam(req, "id").ToInt()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	content, err := router.env.LoadRawContent(id)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(content.Build())
}
