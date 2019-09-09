package repository

import "gitlab.com/ftchinese/content-api/models"

func (env Env) RetrieveRawStory(id string) (models.RawStory, error) {
	var story models.RawStory

	if err := env.db.Get(&story, stmtStory, id); err != nil {
		return models.RawStory{}, err
	}

	story.Sanitize()
	story.SetBilingual()

	return story, nil
}
