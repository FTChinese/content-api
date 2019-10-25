package repository

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

const (
	prefixChannelConfig = "channel_config_"
	keyChannelMap       = "channel_map"
	prefixChannelPage   = "channel_page_"
)

type ChannelEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewChannelEnv(db *sqlx.DB) ChannelEnv {
	return ChannelEnv{
		db:    db,
		cache: cache.New(24*time.Hour, 1*time.Hour),
	}
}

func (env ChannelEnv) TodayFrontPage() (models.FrontPage, error) {
	var teasers = []models.Teaser{}

	if err := env.db.Select(&teasers, stmtFrontPageToday); err != nil {
		return models.FrontPage{}, err
	}

	if len(teasers) > 0 {
		for _, item := range teasers {
			item.Normalize()
		}

		return models.FrontPage{
			Date: chrono.DateFrom(teasers[0].CreatedAt.Time),
			Data: teasers,
		}, nil
	}

	return models.FrontPage{
		Date: chrono.Date{},
		Data: teasers,
	}, nil
}

func (env ChannelEnv) ArchivedFrontPage(date string) (models.ArchivedFrontPage, error) {
	var teasers = []models.Teaser{}

	if err := env.db.Select(&teasers, stmtFrontPageArchive, date); err != nil {
		return models.ArchivedFrontPage{}, err
	}

	if len(teasers) > 0 {
		for _, item := range teasers {
			item.Normalize()
		}
	}

	return models.ArchivedFrontPage{
		Date: date,
		Data: teasers,
	}, nil
}

// ListChannels selects all channels available.
func (env ChannelEnv) ListChannels() ([]models.ChannelSetting, error) {
	var data = make([]models.ChannelSetting, 0)

	if err := env.db.Select(&data, stmtListChannels); err != nil {
		return data, err
	}

	return data, nil
}

// LoadChannelMap loads channels' name to id map from cache
// or from db.
func (env ChannelEnv) LoadChannelMap() (models.ChannelMap, error) {
	if m, ok := env.loadCachedChannelMap(); ok {
		return m, nil
	}

	items, err := env.ListChannels()
	if err != nil {
		return models.ChannelMap{}, err
	}

	m := models.NewChannelMap(items)

	env.CacheChannelMap(m)

	return m, nil
}

func (env ChannelEnv) loadCachedChannelMap() (models.ChannelMap, bool) {
	x, found := env.cache.Get(keyChannelMap)
	if !found {
		return models.ChannelMap{}, false
	}

	if m, ok := x.(models.ChannelMap); ok {
		return m, true
	}

	return models.ChannelMap{}, false
}

func (env ChannelEnv) CacheChannelMap(m models.ChannelMap) {
	env.cache.Set(keyChannelMap, m, cache.DefaultExpiration)
}

// LoadChannel loads a single entry of channel configuration data.
func (env ChannelEnv) LoadChannel(name string) (models.ChannelSetting, error) {
	if c, ok := env.loadCachedChannel(name); ok {
		return c, nil
	}

	c, err := env.retrieveChannel(name)
	if err != nil {
		return c, err
	}

	env.cacheChannel(name, c)

	return c, nil
}

func (env ChannelEnv) loadCachedChannel(name string) (models.ChannelSetting, bool) {
	x, found := env.cache.Get(prefixChannelConfig + name)
	if !found {
		return models.ChannelSetting{}, false
	}

	if c, ok := x.(models.ChannelSetting); ok {
		return c, true
	}

	return models.ChannelSetting{}, false
}

// retrieveChannel retrieves the configuration data of the
// specified channel.
func (env ChannelEnv) retrieveChannel(pathName string) (models.ChannelSetting, error) {
	var c models.ChannelSetting

	err := env.db.Get(&c, stmtChannelSetting, pathName)

	if err != nil {
		return c, err
	}

	return c, nil
}

func (env ChannelEnv) cacheChannel(name string, c models.ChannelSetting) {
	env.cache.Set(prefixChannelConfig+name, c, cache.DefaultExpiration)
}
