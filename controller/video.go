package controller

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
)

type VideoRouter struct {
	env repository.Env
}

func NewVideoRouter(db *sqlx.DB) VideoRouter {
	return VideoRouter{env: repository.NewEnv(db)}
}

func (router VideoRouter) Article(w http.ResponseWriter, req *http.Request) {

}
