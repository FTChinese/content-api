package repository

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	models2 "gitlab.com/ftchinese/content-api/internal/pkg"
	"go.uber.org/zap"
	"time"
)

const (
	prefixChannelConfig = "channel_config_"
	keyChannelMap       = "channel_map"
	prefixChannelPage   = "channel_page_"
)

type ChannelEnv struct {
	db     *sqlx.DB
	cache  *cache.Cache
	logger *zap.Logger
}

func NewChannelEnv(db *sqlx.DB, logger *zap.Logger) ChannelEnv {
	return ChannelEnv{
		db:    db,
		cache: cache.New(24*time.Hour, 1*time.Hour),
	}
}

func (env ChannelEnv) TodayFrontPage() (models2.FrontPage, error) {
	var stories = make([]models2.RawStory, 0)

	if err := env.db.Select(&stories, stmtFrontPageToday); err != nil {
		return models2.FrontPage{}, err
	}

	var teasers = make([]models2.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
	}

	return models2.FrontPage{
		Date: chrono.DateFrom(stories[0].CreatedAt.Time),
		Data: teasers,
	}, nil
}

func (env ChannelEnv) ArchivedFrontPage(date string) (models2.ArchivedFrontPage, error) {
	var stories = make([]models2.RawStory, 0)

	if err := env.db.Select(&stories, stmtFrontPageArchive, date); err != nil {
		return models2.ArchivedFrontPage{}, err
	}

	var teasers = make([]models2.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
	}

	return models2.ArchivedFrontPage{
		Date: date,
		Data: teasers,
	}, nil
}

// ListChannels selects all channels available.
func (env ChannelEnv) ListChannels() ([]models2.ChannelSetting, error) {
	var data = make([]models2.ChannelSetting, 0)

	if err := env.db.Select(&data, stmtListChannels); err != nil {
		return data, err
	}

	return data, nil
}

// LoadChannelMap loads channels' name to id map from cache
// or from db.
func (env ChannelEnv) LoadChannelMap() (models2.ChannelMap, error) {
	if m, ok := env.loadCachedChannelMap(); ok {
		return m, nil
	}

	items, err := env.ListChannels()
	if err != nil {
		return models2.ChannelMap{}, err
	}

	m := models2.NewChannelMap(items)

	env.CacheChannelMap(m)

	return m, nil
}

func (env ChannelEnv) loadCachedChannelMap() (models2.ChannelMap, bool) {
	x, found := env.cache.Get(keyChannelMap)
	if !found {
		return models2.ChannelMap{}, false
	}

	if m, ok := x.(models2.ChannelMap); ok {
		return m, true
	}

	return models2.ChannelMap{}, false
}

func (env ChannelEnv) CacheChannelMap(m models2.ChannelMap) {
	env.cache.Set(keyChannelMap, m, cache.DefaultExpiration)
}

// LoadChannel loads a single entry of channel configuration data.
func (env ChannelEnv) LoadChannel(name string) (models2.ChannelSetting, error) {
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

func (env ChannelEnv) loadCachedChannel(name string) (models2.ChannelSetting, bool) {
	x, found := env.cache.Get(prefixChannelConfig + name)
	if !found {
		return models2.ChannelSetting{}, false
	}

	if c, ok := x.(models2.ChannelSetting); ok {
		return c, true
	}

	return models2.ChannelSetting{}, false
}

// retrieveChannel retrieves the configuration data of the
// specified channel.
func (env ChannelEnv) retrieveChannel(pathName string) (models2.ChannelSetting, error) {
	var c models2.ChannelSetting

	err := env.db.Get(&c, stmtChannelSetting, pathName)

	if err != nil {
		return c, err
	}

	return c, nil
}

// RetrieveTeasers selects all article summary belong to a
// channel.
func (env ChannelEnv) RetrieveTeasers(channelID int64, p gorest.Pagination) ([]models2.Teaser, error) {
	var raws = make([]models2.RawContentBase, 0)
	var t = make([]models2.Teaser, 0)

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

func (env ChannelEnv) cacheChannel(name string, c models2.ChannelSetting) {
	env.cache.Set(prefixChannelConfig+name, c, cache.DefaultExpiration)
}

type chConfigResult struct {
	success models2.ChannelSetting
	err     error
}

type teasersResult struct {
	success []models2.Teaser
	err     error
}

// LoadPage loads a channel's details by retrieving data concurrently.
func (env ChannelEnv) LoadPage(channelName string, p gorest.Pagination) (models2.ChannelPage, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	chMap, err := env.LoadChannelMap()
	if err != nil {
		log.Error(err)
		return models2.ChannelPage{}, err
	}

	chID, ok := chMap[channelName]
	if !ok {
		return models2.ChannelPage{}, err
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
		return models2.ChannelPage{}, cResult.err
	}

	if tResult.err != nil {
		log.Error(tResult.err)
		return models2.ChannelPage{}, tResult.err
	}

	return models2.ChannelPage{
		ChannelSetting: cResult.success,
		Data:           tResult.success,
	}, nil
}
