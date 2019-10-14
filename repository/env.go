package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

type Env struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewEnv(db *sqlx.DB) Env {
	return Env{
		db: db,
		// Default expiration 24 hours, and purges the expired items every hour.
		cache: cache.New(24*time.Hour, 1*time.Hour),
	}
}

var logger = logrus.WithField("package", "content-api.repository")
