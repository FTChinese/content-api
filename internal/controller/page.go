package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"gitlab.com/ftchinese/content-api/internal/repository"
	"gitlab.com/ftchinese/content-api/pkg/xhttp"
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
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(frontPage)
}

func (router PageRouter) ArchivedFrontPage(w http.ResponseWriter, req *http.Request) {
	date, err := xhttp.GetURLParam(req, "date").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	frontPage, err := router.env.ArchivedFrontPage(date)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(frontPage)
}

func (router PageRouter) ChannelList(w http.ResponseWriter, req *http.Request) {
	chs, err := router.env.ListChannels()
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	router.env.CacheChannelMap(pkg.NewChannelMap(chs))

	_ = render.New(w).OK(chs)
}

func (router PageRouter) ChannelData(w http.ResponseWriter, req *http.Request) {
	name, err := xhttp.GetURLParam(req, "name").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	p := gorest.GetPagination(req)
	data, err := router.env.LoadPage(name, p)

	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(data)
}

func (router PageRouter) InspectChannelMap(w http.ResponseWriter, req *http.Request) {
	m, err := router.env.LoadChannelMap()
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(m)
}
