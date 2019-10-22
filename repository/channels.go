package repository

import (
	"github.com/FTChinese/go-rest/chrono"
	"gitlab.com/ftchinese/content-api/models"
)

func (env Env) TodayFrontPage() (models.FrontPage, error) {
	var teasers = []models.Teaser{}

	if err := env.db.Select(&teasers, stmtFrontPageToday); err != nil {
		return models.FrontPage{}, err
	}

	if len(teasers) > 0 {
		for _, item := range teasers {
			item.Normalize()
		}

		return models.FrontPage{
			Date: chrono.DateFrom(teasers[0].CreatedAt.Time),
			Data: teasers,
		}, nil
	}

	return models.FrontPage{
		Date: chrono.Date{},
		Data: teasers,
	}, nil
}

func (env Env) ArchivedFrontPage(date string) (models.ArchivedFrontPage, error) {
	var teasers = []models.Teaser{}

	if err := env.db.Select(&teasers, stmtFrontPageArchive, date); err != nil {
		return models.ArchivedFrontPage{}, err
	}

	if len(teasers) > 0 {
		for _, item := range teasers {
			item.Normalize()
		}
	}

	return models.ArchivedFrontPage{
		Date: date,
		Data: teasers,
	}, nil
}
