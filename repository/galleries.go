package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

type imageResult struct {
	success []models.GalleryItem
	err     error
}

type galleryResult struct {
	success models.Gallery
	err     error
}

type GalleryEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewGalleryEnv(db *sqlx.DB) GalleryEnv {
	return GalleryEnv{
		db:    db,
		cache: cache.New(30*time.Minute, 30*time.Minute),
	}
}

func (env GalleryEnv) LoadGallery(id string) (models.Gallery, error) {
	log := logger.WithField("trace", "GalleryEnv.LoadGallery")

	if g, ok := env.getCachedGallery(id); ok {
		log.Infof("Loaded gallery %s from cache", id)
		return g, nil
	}

	g, err := env.RetrieveGallery(id)
	if err != nil {
		log.Error(err)
		return models.Gallery{}, err
	}

	env.cacheGallery(g)

	log.Infof("Loaded gallery %s from db", id)

	return g, nil
}

func (env GalleryEnv) retrieveGalleryImages(id string) ([]models.GalleryItem, error) {
	var items = make([]models.GalleryItem, 0)

	if err := env.db.Select(&items, stmtGalleryImages, id); err != nil {
		return []models.GalleryItem{}, err
	}

	return items, nil
}

func (env GalleryEnv) retrieveGalleryBody(id string) (models.Gallery, error) {
	var data models.RawGallery

	if err := env.db.Get(&data, stmtGallery, id); err != nil {
		return models.Gallery{}, err
	}

	return data.Build(), nil
}

func (env GalleryEnv) RetrieveGallery(id string) (models.Gallery, error) {
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
		return models.Gallery{}, iResult.err
	}

	if gResult.err != nil {
		return models.Gallery{}, gResult.err
	}

	gResult.success.Items = iResult.success

	return gResult.success, nil
}

func (env GalleryEnv) cacheGallery(g models.Gallery) {
	logger.WithField("trace", "GalleryEnv.cacheGallery").Infof("Caching new gallery %s", g.ID)

	env.cache.Set(g.ID, g, cache.DefaultExpiration)
}

func (env GalleryEnv) getCachedGallery(id string) (models.Gallery, bool) {
	x, found := env.cache.Get(id)
	if !found {
		return models.Gallery{}, false
	}

	if g, ok := x.(models.Gallery); ok {
		return g, true
	}

	return models.Gallery{}, false
}
