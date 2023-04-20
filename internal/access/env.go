package access

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
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

// Load tries to load an access token from cache first, then
// retrieve from db if not found in cache.
func (env Env) Load(token string) (OAuth, error) {
	if acc, ok := env.loadCachedToken(token); ok {
		return acc, nil
	}

	acc, err := env.retrieveFromDB(token)
	if err != nil {
		return acc, err
	}

	env.CacheToken(token, acc)

	return acc, nil
}

func (env Env) loadCachedToken(token string) (OAuth, bool) {
	x, found := env.cache.Get(token)
	if !found {
		return OAuth{}, false
	}

	if access, ok := x.(OAuth); ok {
		return access, true
	}

	return OAuth{}, false
}

func (env Env) retrieveFromDB(token string) (OAuth, error) {
	var access OAuth

	if err := env.db.Get(&access, stmtOAuth, token); err != nil {
		return access, err
	}

	return access, nil
}

func (env Env) CacheToken(token string, access OAuth) {
	env.cache.Set(token, access, cache.DefaultExpiration)
}
