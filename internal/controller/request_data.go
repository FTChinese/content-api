package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = logrus.WithField("package", "controller")

// GetURLParam gets a url parameter.
func GetURLParam(req *http.Request, key string) gorest.Param {
	v := chi.URLParam(req, key)

	return gorest.NewParam(key, v)
}
