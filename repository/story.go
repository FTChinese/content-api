package repository

import "gitlab.com/ftchinese/content-api/models"

type rawStoryResult struct {
	success models.RawStory
	err     error
}

type relatedStoryResult struct {
	success []models.ArticleMeta
	err     error
}

func (env ContentEnv) RetrieveRawStory(id string) (models.RawStory, error) {
	var story models.RawStory

	if err := env.db.Get(&story, stmtStory, id); err != nil {
		logger.WithField("trace", "ContentEnv.RetrieveRawStory").Error(err)
		return models.RawStory{}, err
	}

	story.Related = []models.ArticleMeta{} // So that it won't output null.
	story.Sanitize()
	story.Normalize()

	return story, nil
}

func (env ContentEnv) RelatedStories(id string) ([]models.ArticleMeta, error) {
	var stories []models.RawContentBase

	if err := env.db.Select(&stories, stmtRelatedStory, id); err != nil {
		logger.WithField("trace", "ContentEnv.RelatedStories").Error(err)

		return []models.ArticleMeta{}, err
	}

	var items = make([]models.ArticleMeta, 0)

	for _, v := range stories {
		items = append(items, v.ArticleMeta())
	}

	return items, nil
}

func (env ContentEnv) RawStory(id string) (models.RawStory, error) {
	storyChan := make(chan rawStoryResult)
	relatedChan := make(chan relatedStoryResult)

	go func() {
		story, err := env.RetrieveRawStory(id)
		storyChan <- rawStoryResult{
			success: story,
			err:     err,
		}
	}()

	go func() {
		related, err := env.RelatedStories(id)
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
