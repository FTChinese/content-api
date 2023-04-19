package repository

import (
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"go.uber.org/zap"
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
		db:     db,
		cache:  cache.New(24*time.Hour, 1*time.Hour),
		logger: logger,
	}
}

func (env ChannelEnv) TodayFrontPage() (pkg.FrontPage, error) {
	var stories = make([]pkg.RawStory, 0)

	if err := env.db.Select(&stories, pkg.StmtFrontPageToday); err != nil {
		return pkg.FrontPage{}, err
	}

	var teasers = make([]pkg.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
	}

	return pkg.FrontPage{
		Date: stories[0].CreatedAt,
		Data: teasers,
	}, nil
}

func (env ChannelEnv) ArchivedFrontPage(date string) (pkg.ArchivedFrontPage, error) {
	var stories = make([]pkg.RawStory, 0)

	if err := env.db.Select(&stories, pkg.StmtFrontPageArchive, date); err != nil {
		return pkg.ArchivedFrontPage{}, err
	}

	var teasers = make([]pkg.Teaser, 0)

	for _, item := range stories {
		teasers = append(teasers, item.Teaser())
	}

	return pkg.ArchivedFrontPage{
		Date: date,
		Data: teasers,
	}, nil
}

// ListChannels selects all channels available.
func (env ChannelEnv) ListChannels() ([]pkg.ChannelSetting, error) {
	var data = make([]pkg.ChannelSetting, 0)

	if err := env.db.Select(&data, stmtListChannels); err != nil {
		return data, err
	}

	return data, nil
}

// LoadChannelMap loads channels' name to id map from cache
// or from db.
func (env ChannelEnv) LoadChannelMap() (pkg.ChannelMap, error) {
	if m, ok := env.loadCachedChannelMap(); ok {
		return m, nil
	}

	items, err := env.ListChannels()
	if err != nil {
		return pkg.ChannelMap{}, err
	}

	m := pkg.NewChannelMap(items)

	env.CacheChannelMap(m)

	return m, nil
}

func (env ChannelEnv) loadCachedChannelMap() (pkg.ChannelMap, bool) {
	x, found := env.cache.Get(keyChannelMap)
	if !found {
		return pkg.ChannelMap{}, false
	}

	if m, ok := x.(pkg.ChannelMap); ok {
		return m, true
	}

	return pkg.ChannelMap{}, false
}

func (env ChannelEnv) CacheChannelMap(m pkg.ChannelMap) {
	env.cache.Set(keyChannelMap, m, cache.DefaultExpiration)
}

// LoadChannel loads a single entry of channel configuration data.
func (env ChannelEnv) LoadChannel(name string) (pkg.ChannelSetting, error) {
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

func (env ChannelEnv) loadCachedChannel(name string) (pkg.ChannelSetting, bool) {
	x, found := env.cache.Get(prefixChannelConfig + name)
	if !found {
		return pkg.ChannelSetting{}, false
	}

	if c, ok := x.(pkg.ChannelSetting); ok {
		return c, true
	}

	return pkg.ChannelSetting{}, false
}

// retrieveChannel retrieves the configuration data of the
// specified channel.
func (env ChannelEnv) retrieveChannel(pathName string) (pkg.ChannelSetting, error) {
	var c pkg.ChannelSetting

	err := env.db.Get(&c, stmtChannelSetting, pathName)

	if err != nil {
		return c, err
	}

	return c, nil
}

// RetrieveTeasers selects all article summary belong to a
// channel.
func (env ChannelEnv) RetrieveTeasers(channelID int64, p gorest.Pagination) ([]pkg.Teaser, error) {
	var raws = make([]pkg.RawContentBase, 0)
	var t = make([]pkg.Teaser, 0)

	err := env.db.Select(
		&raws,
		pkg.StmtStoryTeaser,
		channelID,
		p.Limit,
		p.Offset(),
	)
	if err != nil {
		return t, err
	}

	for _, v := range raws {
		t = append(t, v.Teaser())
	}

	return t, err
}

func (env ChannelEnv) cacheChannel(name string, c pkg.ChannelSetting) {
	env.cache.Set(prefixChannelConfig+name, c, cache.DefaultExpiration)
}

type chConfigResult struct {
	success pkg.ChannelSetting
	err     error
}

type teasersResult struct {
	success []pkg.Teaser
	err     error
}

// LoadPage loads a channel's details by retrieving data concurrently.
func (env ChannelEnv) LoadPage(channelName string, p gorest.Pagination) (pkg.ChannelPage, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	chMap, err := env.LoadChannelMap()
	if err != nil {
		log.Error(err)
		return pkg.ChannelPage{}, err
	}

	chID, ok := chMap[channelName]
	if !ok {
		return pkg.ChannelPage{}, err
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
		return pkg.ChannelPage{}, cResult.err
	}

	if tResult.err != nil {
		log.Error(tResult.err)
		return pkg.ChannelPage{}, tResult.err
	}

	return pkg.ChannelPage{
		ChannelSetting: cResult.success,
		Data:           tResult.success,
	}, nil
}
