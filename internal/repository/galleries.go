package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"go.uber.org/zap"
)

type imageResult struct {
	success []pkg.GalleryItem
	err     error
}

type galleryResult struct {
	success pkg.Gallery
	err     error
}

type GalleryEnv struct {
	db     *sqlx.DB
	cache  *cache.Cache
	logger *zap.Logger
}

func NewGalleryEnv(db *sqlx.DB, logger *zap.Logger) GalleryEnv {
	return GalleryEnv{
		db:     db,
		cache:  cache.New(30*time.Minute, 30*time.Minute),
		logger: logger,
	}
}

func (env GalleryEnv) LoadGallery(id string) (pkg.Gallery, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	if g, ok := env.getCachedGallery(id); ok {
		log.Infof("Loaded gallery %s from cache", id)
		return g, nil
	}

	g, err := env.RetrieveGallery(id)
	if err != nil {
		log.Error(err)
		return pkg.Gallery{}, err
	}

	env.cacheGallery(g)

	log.Infof("Loaded gallery %s from db", id)

	return g, nil
}

func (env GalleryEnv) retrieveGalleryImages(id string) ([]pkg.GalleryItem, error) {
	var items = make([]pkg.GalleryItem, 0)

	if err := env.db.Select(&items, pkg.StmtGalleryItem, id); err != nil {
		return []pkg.GalleryItem{}, err
	}

	return items, nil
}

func (env GalleryEnv) retrieveGalleryBody(id string) (pkg.Gallery, error) {
	var data pkg.RawGallery

	if err := env.db.Get(&data, pkg.StmtGallery, id); err != nil {
		return pkg.Gallery{}, err
	}

	return data.Build(), nil
}

func (env GalleryEnv) RetrieveGallery(id string) (pkg.Gallery, error) {
	imageChan := make(chan imageResult)
	bodyChan := make(chan galleryResult)

	go func() {
		images, err := env.retrieveGalleryImages(id)
		imageChan <- imageResult{
			success: images,
			err:     err,
		}
	}()

	go func() {
		body, err := env.retrieveGalleryBody(id)
		bodyChan <- galleryResult{
			success: body,
			err:     err,
		}
	}()

	iResult, gResult := <-imageChan, <-bodyChan

	if iResult.err != nil {
		return pkg.Gallery{}, iResult.err
	}

	if gResult.err != nil {
		return pkg.Gallery{}, gResult.err
	}

	gResult.success.Items = iResult.success

	return gResult.success, nil
}

func (env GalleryEnv) cacheGallery(g pkg.Gallery) {
	env.cache.Set(g.ID, g, cache.DefaultExpiration)
}

func (env GalleryEnv) getCachedGallery(id string) (pkg.Gallery, bool) {
	x, found := env.cache.Get(id)
	if !found {
		return pkg.Gallery{}, false
	}

	if g, ok := x.(pkg.Gallery); ok {
		return g, true
	}

	return pkg.Gallery{}, false
}
