package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

type VideoEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewVideoEnv(db *sqlx.DB) VideoEnv {
	return VideoEnv{
		db:    db,
		cache: cache.New(30*time.Minute, 30*time.Minute),
	}
}

func (env VideoEnv) Load(id int64) (models.Video, error) {
	log := logger.WithField("trace", "VideoEnv.Load")

	if v, ok := env.getCachedVideo(id); ok {
		log.Infof("Loaded video %d from cache", id)
		return v, nil
	}

	v, err := env.RetrieveVideo(id)
	if err != nil {
		log.Error(err)
		return models.Video{}, err
	}

	env.cacheVideo(v)

	log.Infof("Loaded video %d from db", id)

	return v, nil
}

func (env VideoEnv) RetrieveVideo(id int64) (models.Video, error) {
	var video models.RawVideo

	if err := env.db.Get(&video, stmtVideo, id); err != nil {
		return models.Video{}, err
	}

	return video.Build(), nil
}

func (env VideoEnv) cacheVideo(v models.Video) {
	logger.WithField("trace", "VideoEnv.cacheVideo").Infof("Caching video %d", v.ID)

	env.cache.Set(string(v.ID), v, cache.DefaultExpiration)
}

func (env VideoEnv) getCachedVideo(id int64) (models.Video, bool) {
	x, found := env.cache.Get(string(id))

	if !found {
		return models.Video{}, false
	}

	if v, ok := x.(models.Video); ok {
		return v, true
	}

	return models.Video{}, false
}
