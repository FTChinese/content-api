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
		KeyWords:    null.StringFrom("英语电台,音乐"),
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
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	//"english_bilingual-reading": {
	//    ID:          0,
	//    ParentID:    0,
	//    KeyName:     "english_bilingual-reading",
	//    Name:        "bilingual-reading",
	//    Title:       "双语阅读",
	//    Description: null.String{},
	//    KeyWords:    null.String{},
	//    CreatedAt:   chrono.Time{},
	//    UpdatedAt:   chrono.Time{},
	//},
	"ft-radio": {
		ID:          0,
		ParentID:    0,
		KeyName:     "ft-radio",
		Name:        "ft-radio",
		Title:       "英语电台",
		Description: null.String{},
		KeyWords:    null.StringFrom("英语电台"),
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
	//"english_speed-reading": {
	//    ID:          0,
	//    ParentID:    0,
	//    KeyName:     "english_speed-reading",
	//    Name:        "speed-reading",
	//    Title:       "FT英语速读",
	//    Description: null.String{},
	//    KeyWords:    null.StringFrom("速读"),
	//    CreatedAt:   chrono.Time{},
	//    UpdatedAt:   chrono.Time{},
	//},
	"daily-word": {
		ID:          0,
		ParentID:    0,
		KeyName:     "daily-word",
		Name:        "daily-word",
		Title:       "每日一词",
		Description: null.String{},
		KeyWords:    null.StringFrom("每日一词"),
		CreatedAt:   chrono.Time{},
		UpdatedAt:   chrono.Time{},
	},
}

var SpeedReadingChannel = models.ChannelSetting{
	ID:          0,
	ParentID:    0,
	KeyName:     "speed-reading",
	Name:        "speed-reading",
	Title:       "FT英语速读",
	Description: null.String{},
	KeyWords:    null.StringFrom("速读"),
	CreatedAt:   chrono.Time{},
	UpdatedAt:   chrono.Time{},
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

// RetrieveChannelTeasers loads the channel article list.
func (env InteractiveEnv) RetrieveChannelTeasers(keyWords string, p gorest.Pagination) ([]models.Teaser, error) {
	var teasers = make([]models.Teaser, 0)

	err := env.db.Select(&teasers, stmtAudioTeasers,
		keyWords,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		return teasers, err
	}

	for i := range teasers {
		teasers[i].Normalize()
	}

	return teasers, nil
}

// RetrieveAudioArticle loads article of interactive.
func (env InteractiveEnv) RetrieveAudioArticle(id int64) (models.AudioArticle, error) {
	var a models.AudioArticle

	if err := env.db.Get(&a, stmtAudioArticle, id); err != nil {
		return a, err
	}

	a.Normalize()

	return a, nil
}

func (env InteractiveEnv) RetrieveSpeedReadingTeasers(p gorest.Pagination) ([]models.Teaser, error) {
	var teasers = make([]models.Teaser, 0)

	err := env.db.Select(&teasers,
		stmtSpeedReadingTeasers,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		return teasers, err
	}

	for i, _ := range teasers {
		teasers[i].Normalize()
	}

	return teasers, nil
}

func (env InteractiveEnv) RetrieveSpeedReading(id int64) (models.SpeedReading, error) {
	var s models.SpeedReading

	if err := env.db.Get(&s, stmtSpeedReadingContent, id); err != nil {
		return s, err
	}

	s.Normalize()

	return s, nil
}
