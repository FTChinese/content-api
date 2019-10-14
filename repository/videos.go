package repository

import "gitlab.com/ftchinese/content-api/models"

func (env Env) RetrieveVideo(id int64) (models.Video, error) {
	var video models.Video

	if err := env.db.Get(&video, stmtVideo, id); err != nil {
		return video, err
	}

	return video, nil
}
