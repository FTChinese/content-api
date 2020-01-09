package repository

import (
	gorest "github.com/FTChinese/go-rest"
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
	var stories = make([]models.RawStory, 0)

	if err := env.db.Select(&stories, stmtFrontPageToday); err != nil {
		return models.FrontPage{}, err
	}

	var teasers = make([]models.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
	}

	return models.FrontPage{
		Date: chrono.DateFrom(stories[0].CreatedAt.Time),
		Data: teasers,
	}, nil
}

func (env ChannelEnv) ArchivedFrontPage(date string) (models.ArchivedFrontPage, error) {
	var stories = make([]models.RawStory, 0)

	if err := env.db.Select(&stories, stmtFrontPageArchive, date); err != nil {
		return models.ArchivedFrontPage{}, err
	}

	var teasers = make([]models.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
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
		logger.WithField("trace", "ChannelEnv.ListChannels")

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

// RetrieveTeasers selects all article summary belong to a
// channel.
func (env ChannelEnv) RetrieveTeasers(channelID int64, p gorest.Pagination) ([]models.Teaser, error) {
	var raws = make([]models.RawContentBase, 0)
	var t = make([]models.Teaser, 0)

	err := env.db.Select(&raws, stmtStoryTeaser,
		channelID, p.Limit, p.Offset(),
	)
	if err != nil {
		return t, err
	}

	for _, v := range raws {
		t = append(t, v.Teaser())
	}

	return t, err
}

func (env ChannelEnv) cacheChannel(name string, c models.ChannelSetting) {
	env.cache.Set(prefixChannelConfig+name, c, cache.DefaultExpiration)
}

type chConfigResult struct {
	success models.ChannelSetting
	err     error
}

type teasersResult struct {
	success []models.Teaser
	err     error
}

// LoadPage loads a channel's details by retrieving data concurrently.
func (env ChannelEnv) LoadPage(channelName string, p gorest.Pagination) (models.ChannelPage, error) {
	log := logger.WithField("trace", "ChannelEnv.LoadPage")

	chMap, err := env.LoadChannelMap()
	if err != nil {
		log.Error(err)
		return models.ChannelPage{}, err
	}

	chID, ok := chMap[channelName]
	if !ok {
		return models.ChannelPage{}, err
	}

	configChan := make(chan chConfigResult)
	teasersChan := make(chan teasersResult)

	go func() {
		chConfig, err := env.LoadChannel(channelName)

		configChan <- chConfigResult{
			success: chConfig,
			err:     err,
		}
	}()

	go func() {
		teasers, err := env.RetrieveTeasers(chID, p)

		teasersChan <- teasersResult{
			success: teasers,
			err:     err,
		}
	}()

	cResult, tResult := <-configChan, <-teasersChan

	if cResult.err != nil {
		log.Error(cResult.err)
		return models.ChannelPage{}, cResult.err
	}

	if tResult.err != nil {
		log.Error(tResult.err)
		return models.ChannelPage{}, tResult.err
	}

	return models.ChannelPage{
		ChannelSetting: cResult.success,
		Data:           tResult.success,
	}, nil
}
