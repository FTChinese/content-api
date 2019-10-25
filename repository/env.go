package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

type ContentEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewContentEnv(db *sqlx.DB) ContentEnv {
	return ContentEnv{
		db: db,
		// Default expiration 24 hours, and purges the expired items every hour.
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

var logger = logrus.WithField("package", "content-api.repository")
