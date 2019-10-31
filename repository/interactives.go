package repository

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

var interactiveMap = map[string]models.ChannelSetting{
	"mle": {
		ID:          0,
		ParentID:    0,
		KeyName:     "mle",
		Name:        "mle",
		Title:       "麦可林学英语",
		Description: null.String{},
		KeyWords:    null.StringFrom("麦可林学英语"),
		ContentKind: models.ContentKindAudio,
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
		ContentKind: models.ContentKindAudio,
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
		ContentKind: models.ContentKindAudio,
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
		ContentKind: models.ContentKindAudio,
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
		ContentKind: models.ContentKindAudio,
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
		ContentKind: models.ContentKindSpeedReading,
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
		ContentKind: models.ContentKindAudio,
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
}

func GetAudioChannelConfig(name string) (models.ChannelSetting, bool) {
	config, ok := interactiveMap[name]
	return config, ok
}

type InteractiveEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewInteractiveEnv(db *sqlx.DB) InteractiveEnv {
	return InteractiveEnv{
		db:    db,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// RetrieveTeasers loads a list article summary.
func (env InteractiveEnv) RetrieveTeasers(keyWords string, p gorest.Pagination) ([]models.RawInteractive, error) {
	var teasers = make([]models.RawInteractive, 0)

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

func (env InteractiveEnv) RetrieveRawContent(id int64) (models.RawInteractive, error) {
	var c models.RawInteractive

	if err := env.db.Get(&c, stmtInteractiveContent, id); err != nil {
		return c, err
	}

	return c, nil
}
