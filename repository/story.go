package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"gitlab.com/ftchinese/content-api/models"
	"time"
)

type rawStoryResult struct {
	success models.RawStory
	err     error
}

type relatedStoryResult struct {
	success []models.ArticleMeta
	err     error
}

type StoryEnv struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewStoryEnv(db *sqlx.DB) StoryEnv {
	return StoryEnv{
		db:    db,
		cache: cache.New(30*time.Minute, 30*time.Minute),
	}
}

// LoadRawStory loads raw story from cache first, and from
// db if not found and cache it.
func (env StoryEnv) LoadRawStory(id string) (models.RawStory, error) {
	log := logger.WithField("trace", "StoryEnv.LoadRawStory")

	if story, ok := env.getCachedRawStory(id); ok {
		log.Infof("Loaded raw story %s from cache", id)
		return story, nil
	}

	story, err := env.RetrieveRawStory(id)
	if err != nil {
		log.Error(err)
		return models.RawStory{}, err
	}

	env.cacheRawStory(story)

	log.Infof("Loaded raw story %s from db", id)
	return story, nil
}

func (env StoryEnv) retrieveRawStory(id string) (models.RawStory, error) {
	var story models.RawStory

	if err := env.db.Get(&story, stmtStory, id); err != nil {
		logger.WithField("trace", "StoryEnv.retrieveRawStory").Errorf("Story %s, %s", id, err)
		return models.RawStory{}, err
	}

	story.Related = []models.ArticleMeta{} // So that it won't output null.
	story.Sanitize()
	story.Normalize()

	return story, nil
}

func (env StoryEnv) retrieveRelatedStories(id string) ([]models.ArticleMeta, error) {
	var stories []models.RawContentBase

	if err := env.db.Select(&stories, stmtRelatedStory, id); err != nil {
		logger.WithField("trace", "ContentEnv.retrieveRelatedStories").Error(err)

		return []models.ArticleMeta{}, err
	}

	var items = make([]models.ArticleMeta, 0)

	for _, v := range stories {
		items = append(items, v.ArticleMeta())
	}

	return items, nil
}

// RetrieveRawStory retrieves a story and its related
// articles form DB.
func (env StoryEnv) RetrieveRawStory(id string) (models.RawStory, error) {
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
		return models.RawStory{}, storyResult.err
	}

	// If related stories are not retrieved, leave it empty.
	if relatedResult.err != nil {
		return storyResult.success, nil
	}

	storyResult.success.Related = relatedResult.success

	return storyResult.success, nil
}

func (env StoryEnv) cacheRawStory(raw models.RawStory) {
	logger.WithField("trace", "StoryEnv.cacheRawStory").Infof("Caching raw story %s", raw.ID)

	env.cache.Set(raw.ID, raw, cache.DefaultExpiration)
}

func (env StoryEnv) getCachedRawStory(id string) (models.RawStory, bool) {
	x, found := env.cache.Get(id)
	if !found {
		return models.RawStory{}, false
	}

	if story, ok := x.(models.RawStory); ok {
		return story, true
	}

	return models.RawStory{}, false
}
