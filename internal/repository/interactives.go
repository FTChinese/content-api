package repository

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	models2 "gitlab.com/ftchinese/content-api/internal/pkg"
	"go.uber.org/zap"
	"time"
)

var interactiveMap = map[string]models2.ChannelSetting{
	"mle": {
		ID:          0,
		ParentID:    0,
		KeyName:     "mle",
		Name:        "mle",
		Title:       "麦可林学英语",
		Description: null.String{},
		KeyWords:    null.StringFrom("麦可林学英语"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"life-of-a-song": {
		ID:          0,
		ParentID:    0,
		KeyName:     "life-of-a-song",
		Name:        "life-of-a-song",
		Title:       "音乐之生",
		Description: null.String{},
		KeyWords:    null.StringFrom("音乐"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"boomear": {
		ID:          0,
		ParentID:    0,
		KeyName:     "boomear",
		Name:        "boomear",
		Title:       "BoomEar艺术播客",
		Description: null.String{},
		KeyWords:    null.StringFrom("BoomEar艺术播客"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"yibohaoshu": {
		ID:          0,
		ParentID:    0,
		KeyName:     "yibohaoshu",
		Name:        "yibohaoshu",
		Title:       "一波好书",
		Description: null.String{},
		KeyWords:    null.StringFrom("一波好书"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"ft-radio": {
		ID:          0,
		ParentID:    0,
		KeyName:     "ft-radio",
		Name:        "ft-radio",
		Title:       "英语电台",
		Description: null.String{},
		KeyWords:    null.StringFrom("英语电台"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"speed-reading": {
		ID:          0,
		ParentID:    0,
		KeyName:     "speed-reading",
		Name:        "speed-reading",
		Title:       "FT英语速读",
		Description: null.String{},
		KeyWords:    null.StringFrom("速读"),
		ContentKind: models2.ContentKindSpeedReading,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	"daily-word": {
		ID:          0,
		ParentID:    0,
		KeyName:     "daily-word",
		Name:        "daily-word",
		Title:       "每日一词",
		Description: null.String{},
		KeyWords:    null.StringFrom("每日一词"),
		ContentKind: models2.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
}

func GetAudioChannelConfig(name string) (models2.ChannelSetting, bool) {
	config, ok := interactiveMap[name]
	return config, ok
}

type InteractiveEnv struct {
	db     *sqlx.DB
	cache  *cache.Cache
	logger *zap.Logger
}

func NewInteractiveEnv(db *sqlx.DB, logger *zap.Logger) InteractiveEnv {
	return InteractiveEnv{
		db:     db,
		cache:  cache.New(5*time.Minute, 10*time.Minute),
		logger: logger,
	}
}

// RetrieveTeasers loads a list article summary.
func (env InteractiveEnv) RetrieveTeasers(keyWords string, p gorest.Pagination) ([]models2.RawInteractive, error) {
	var teasers = make([]models2.RawInteractive, 0)

	err := env.db.Select(&teasers, stmtInteractiveTeaser,
		keyWords,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		return teasers, err
	}

	return teasers, nil
}

func (env InteractiveEnv) LoadRawContent(id int64) (models2.RawInteractive, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	if i, ok := env.getCachedContent(id); ok {
		log.Infof("Loaded raw interactive article %d from cache", id)
		return i, nil
	}

	i, err := env.RetrieveRawContent(id)
	if err != nil {
		log.Error(err)
		return models2.RawInteractive{}, err
	}

	return i, nil
}

func (env InteractiveEnv) RetrieveRawContent(id int64) (models2.RawInteractive, error) {
	var c models2.RawInteractive

	if err := env.db.Get(&c, stmtInteractiveContent, id); err != nil {
		return c, err
	}

	return c, nil
}

func (env InteractiveEnv) cacheRawContent(r models2.RawInteractive) {
	env.cache.Set(string(r.ID), r, cache.DefaultExpiration)
}

func (env InteractiveEnv) getCachedContent(id int64) (models2.RawInteractive, bool) {
	x, found := env.cache.Get(string(id))
	if !found {
		return models2.RawInteractive{}, false
	}

	if r, ok := x.(models2.RawInteractive); ok {
		return r, true
	}

	return models2.RawInteractive{}, false
}
