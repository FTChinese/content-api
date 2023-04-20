package repository

import (
	"time"

	"github.com/FTchinese/content-api/internal/pkg"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type StoryEnv struct {
	db     *sqlx.DB
	cache  *cache.Cache
	logger *zap.Logger
}

func NewStoryEnv(db *sqlx.DB, logger *zap.Logger) StoryEnv {
	return StoryEnv{
		db:     db,
		cache:  cache.New(30*time.Minute, 30*time.Minute),
		logger: logger,
	}
}

// LoadRawStory loads raw story from cache first, and from
// db if not found and cache it.
func (env StoryEnv) LoadRawStory(id string) (pkg.RawStory, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	if story, ok := env.getCachedRawStory(id); ok {
		log.Infof("Loaded raw story %s from cache", id)
		return story, nil
	}

	story, err := env.RetrieveRawStory(id)
	if err != nil {
		log.Error(err)
		return pkg.RawStory{}, err
	}

	env.cacheRawStory(story)

	log.Infof("Loaded raw story %s from db", id)
	return story, nil
}

func (env StoryEnv) retrieveRawStory(id string) (pkg.RawStory, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	var story pkg.RawStory

	if err := env.db.Get(&story, pkg.StmtStory, id); err != nil {
		log.Error(err)
		return pkg.RawStory{}, err
	}

	story.Related = []pkg.ArticleMeta{} // So that it won't output null.
	story.Sanitize()
	story.Normalize()

	return story, nil
}

func (env StoryEnv) retrieveRelatedStories(id string) ([]pkg.ArticleMeta, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	var stories []pkg.RawContentBase

	if err := env.db.Select(&stories, pkg.StmtRelatedStory, id); err != nil {
		log.Error(err)

		return []pkg.ArticleMeta{}, err
	}

	var items = make([]pkg.ArticleMeta, 0)

	for _, v := range stories {
		items = append(items, v.ArticleMeta())
	}

	return items, nil
}

// RetrieveRawStory retrieves a story and its related
// articles form DB.
func (env StoryEnv) RetrieveRawStory(id string) (pkg.RawStory, error) {
	storyChan := make(chan pkg.AsyncResult[pkg.RawStory])
	relatedChan := make(chan pkg.AsyncResult[[]pkg.ArticleMeta])

	go func() {
		story, err := env.retrieveRawStory(id)
		storyChan <- pkg.AsyncResult[pkg.RawStory]{
			Value: story,
			Err:   err,
		}
	}()

	go func() {
		related, err := env.retrieveRelatedStories(id)
		relatedChan <- pkg.AsyncResult[[]pkg.ArticleMeta]{
			Value: related,
			Err:   err,
		}
	}()

	storyResult, relatedResult := <-storyChan, <-relatedChan

	if storyResult.Err != nil {
		return pkg.RawStory{}, storyResult.Err
	}

	rawStory := storyResult.Value
	rawStory.Related = relatedResult.Value

	return rawStory, nil
}

func (env StoryEnv) cacheRawStory(raw pkg.RawStory) {
	env.cache.Set(raw.ID, raw, cache.DefaultExpiration)
}

func (env StoryEnv) getCachedRawStory(id string) (pkg.RawStory, bool) {
	x, found := env.cache.Get(id)
	if !found {
		return pkg.RawStory{}, false
	}

	if story, ok := x.(pkg.RawStory); ok {
		return story, true
	}

	return pkg.RawStory{}, false
}
