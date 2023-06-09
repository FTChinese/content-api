package repository

import (
	"strconv"
	"time"

	"github.com/FTchinese/content-api/internal/pkg"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type VideoEnv struct {
	db     *sqlx.DB
	cache  *cache.Cache
	logger *zap.Logger
}

func NewVideoEnv(db *sqlx.DB, logger *zap.Logger) VideoEnv {
	return VideoEnv{
		db:     db,
		cache:  cache.New(30*time.Minute, 30*time.Minute),
		logger: logger,
	}
}

func (env VideoEnv) Load(id int64) (pkg.Video, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	if v, ok := env.getCachedVideo(id); ok {
		log.Infof("Loaded video %d from cache", id)
		return v, nil
	}

	v, err := env.RetrieveVideo(id)
	if err != nil {
		log.Error(err)
		return pkg.Video{}, err
	}

	env.cacheVideo(v)

	log.Infof("Loaded video %d from db", id)

	return v, nil
}

func (env VideoEnv) RetrieveVideo(id int64) (pkg.Video, error) {
	var video pkg.RawVideo

	if err := env.db.Get(&video, pkg.StmtVideo, id); err != nil {
		return pkg.Video{}, err
	}

	return video.Build(), nil
}

func (env VideoEnv) cacheVideo(v pkg.Video) {

	env.cache.Set(strconv.FormatInt(v.ID, 10), v, cache.DefaultExpiration)
}

func (env VideoEnv) getCachedVideo(id int64) (pkg.Video, bool) {
	x, found := env.cache.Get(strconv.FormatInt(id, 10))

	if !found {
		return pkg.Video{}, false
	}

	if v, ok := x.(pkg.Video); ok {
		return v, true
	}

	return pkg.Video{}, false
}
