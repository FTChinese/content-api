package repository

import "gitlab.com/ftchinese/content-api/models"

func (env ContentEnv) RetrieveVideo(id int64) (models.Video, error) {
	var video models.RawVideo

	if err := env.db.Get(&video, stmtVideo, id); err != nil {
		return models.Video{}, err
	}

	return models.NewVideo(&video), nil
}
