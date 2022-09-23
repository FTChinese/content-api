package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	models2 "gitlab.com/ftchinese/content-api/internal/pkg"
	"go.uber.org/zap"
	"time"
)

type rawStoryResult struct {
	success models2.RawStory
	err     error
}

type relatedStoryResult struct {
	success []models2.ArticleMeta
	err     error
}

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
func (env StoryEnv) LoadRawStory(id string) (models2.RawStory, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	if story, ok := env.getCachedRawStory(id); ok {
		log.Infof("Loaded raw story %s from cache", id)
		return story, nil
	}

	story, err := env.RetrieveRawStory(id)
	if err != nil {
		log.Error(err)
		return models2.RawStory{}, err
	}

	env.cacheRawStory(story)

	log.Infof("Loaded raw story %s from db", id)
	return story, nil
}

func (env StoryEnv) retrieveRawStory(id string) (models2.RawStory, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	var story models2.RawStory

	if err := env.db.Get(&story, stmtStory, id); err != nil {
		log.Error(err)
		return models2.RawStory{}, err
	}

	story.Related = []models2.ArticleMeta{} // So that it won't output null.
	story.Sanitize()
	story.Normalize()

	return story, nil
}

func (env StoryEnv) retrieveRelatedStories(id string) ([]models2.ArticleMeta, error) {
	defer env.logger.Sync()
	log := env.logger.Sugar()

	var stories []models2.RawContentBase

	if err := env.db.Select(&stories, stmtRelatedStory, id); err != nil {
		log.Error(err)

		return []models2.ArticleMeta{}, err
	}

	var items = make([]models2.ArticleMeta, 0)

	for _, v := range stories {
		items = append(items, v.ArticleMeta())
	}

	return items, nil
}

// RetrieveRawStory retrieves a story and its related
// articles form DB.
func (env StoryEnv) RetrieveRawStory(id string) (models2.RawStory, error) {
	storyChan := make(chan rawStoryResult)
	relatedChan := make(chan relatedStoryResult)

	go func() {
		story, err := env.retrieveRawStory(id)
		storyChan <- rawStoryResult{
			success: story,
			err:     err,
		}
	}()

	go func() {
		related, err := env.retrieveRelatedStories(id)
		relatedChan <- relatedStoryResult{
			success: related,
			err:     err,
		}
	}()

	storyResult, relatedResult := <-storyChan, <-relatedChan

	if storyResult.err != nil {
		return models2.RawStory{}, storyResult.err
	}

	// If related stories are not retrieved, leave it empty.
	if relatedResult.err != nil {
		return storyResult.success, nil
	}

	storyResult.success.Related = relatedResult.success

	return storyResult.success, nil
}

func (env StoryEnv) cacheRawStory(raw models2.RawStory) {
	env.cache.Set(raw.ID, raw, cache.DefaultExpiration)
}

func (env StoryEnv) getCachedRawStory(id string) (models2.RawStory, bool) {
	x, found := env.cache.Get(id)
	if !found {
		return models2.RawStory{}, false
	}

	if story, ok := x.(models2.RawStory); ok {
		return story, true
	}

	return models2.RawStory{}, false
}
