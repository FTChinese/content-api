package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"gitlab.com/ftchinese/content-api/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

type PageRouter struct {
	env repository.ChannelEnv
}

func NewPageRouter(db *sqlx.DB, logger *zap.Logger) PageRouter {
	return PageRouter{
		env: repository.NewChannelEnv(db, logger),
	}
}

func (router PageRouter) TodayFrontPage(w http.ResponseWriter, req *http.Request) {
	frontPage, err := router.env.TodayFrontPage()

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(frontPage))
}

func (router PageRouter) ArchivedFrontPage(w http.ResponseWriter, req *http.Request) {
	date, err := GetURLParam(req, "date").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	frontPage, err := router.env.ArchivedFrontPage(date)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(frontPage))
}

func (router PageRouter) ChannelList(w http.ResponseWriter, req *http.Request) {
	chs, err := router.env.ListChannels()
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	router.env.CacheChannelMap(pkg.NewChannelMap(chs))

	_ = view.Render(w, view.NewResponse().SetBody(chs))
}

func (router PageRouter) ChannelData(w http.ResponseWriter, req *http.Request) {
	name, err := GetURLParam(req, "name").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p := gorest.GetPagination(req)
	data, err := router.env.LoadPage(name, p)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(data))
}

func (router PageRouter) InspectChannelMap(w http.ResponseWriter, req *http.Request) {
	m, err := router.env.LoadChannelMap()
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(m))
}
