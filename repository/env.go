package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Env struct {
	db *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{db: db}
}

var logger = logrus.WithField("package", "content-api.repository")
