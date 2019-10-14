package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

type OAuthEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewOAuthEnv(db *sqlx.DB) OAuthEnv {
	return OAuthEnv{
		db: db,
		// Default expiration 24 hours, and purges the expired items every hour.
		cache: cache.New(24*time.Hour, 1*time.Hour),
	}
}

// Load tries to load an access token from cache first, then
// retrieve from db if not found in cache.
func (env OAuthEnv) Load(token string) (models.OAuthAccess, error) {
	if acc, ok := env.LoadCachedToken(token); ok {
		return acc, nil
	}

	acc, err := env.RetrieveFromDB(token)
	if err != nil {
		return acc, err
	}

	env.CacheToken(token, acc)

	return acc, nil
}

func (env OAuthEnv) RetrieveFromDB(token string) (models.OAuthAccess, error) {
	var access models.OAuthAccess

	if err := env.db.Get(&access, stmtOAuth, token); err != nil {
		return access, err
	}

	return access, nil
}

func (env OAuthEnv) CacheToken(token string, access models.OAuthAccess) {
	env.cache.Set(token, access, cache.DefaultExpiration)
}

func (env OAuthEnv) LoadCachedToken(token string) (models.OAuthAccess, bool) {
	x, found := env.cache.Get(token)
	if !found {
		return models.OAuthAccess{}, false
	}

	if access, ok := x.(models.OAuthAccess); ok {
		return access, true
	}

	return models.OAuthAccess{}, false
}
