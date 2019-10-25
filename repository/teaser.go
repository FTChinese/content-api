package repository

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/content-api/models"
)

type chConfigResult struct {
	success models.ChannelSetting
	err     error
}

type teasersResult struct {
	success []models.Teaser
	err     error
}

// RetrieveTeasers selects all article summary belong to a
// channel.
func (env ChannelEnv) RetrieveTeasers(channelID int64, p gorest.Pagination) ([]models.Teaser, error) {
	var t = make([]models.Teaser, 0)

	err := env.db.Select(&t, stmtChannelContent,
		channelID, p.Limit, p.Offset(),
	)
	if err != nil {
		return t, err
	}

	return t, err
}

func (env ChannelEnv) LoadPage(channelName string, p gorest.Pagination) (models.ChannelPage, error) {
	chMap, err := env.LoadChannelMap()
	if err != nil {
		return models.ChannelPage{}, err
	}

	chID, ok := chMap[channelName]
	if !ok {
		return models.ChannelPage{}, err
	}

	configChan := make(chan chConfigResult)
	teasersChan := make(chan teasersResult)

	go func() {
		chConfig, err := env.LoadChannel(channelName)

		configChan <- chConfigResult{
			success: chConfig,
			err:     err,
		}
	}()

	go func() {
		teasers, err := env.RetrieveTeasers(chID, p)

		teasersChan <- teasersResult{
			success: teasers,
			err:     err,
		}
	}()

	cResult, tResult := <-configChan, <-teasersChan

	if cResult.err != nil {
		return models.ChannelPage{}, cResult.err
	}

	if tResult.err != nil {
		return models.ChannelPage{}, tResult.err
	}

	return models.ChannelPage{
		ChannelSetting: cResult.success,
		Data:           tResult.success,
	}, nil
}
