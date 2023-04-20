package repository

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"gitlab.com/ftchinese/content-api/pkg/collection"
	"gitlab.com/ftchinese/content-api/pkg/db"
	"gitlab.com/ftchinese/content-api/pkg/ids"
	"go.uber.org/zap"
)

type StarRepo struct {
	dbs    db.ReadWriteMyDBs
	logger *zap.Logger
}

func NewStarRepo(dbs db.ReadWriteMyDBs, logger *zap.Logger) StarRepo {
	return StarRepo{
		dbs:    dbs,
		logger: logger,
	}
}

// StarArticle adds an article to user's favorites.
// Duplicate error will be ignored.
func (env StarRepo) StarArticle(params pkg.StarParams) error {

	_, err := env.dbs.Write.NamedExec(pkg.StmtSaveStar, params)

	if err != nil {
		return err
	}

	return nil
}

// UnstarArticle delete an article from user favourites
func (env StarRepo) UnstarArticle(userID, articleID string) error {

	_, err := env.dbs.Delete.Exec(pkg.StmtUnstar, userID, articleID)

	if err != nil {
		return err
	}

	return nil
}

// StmtIsStarring checks if the user is starring an article
func (env StarRepo) IsStarring(userIDs ids.UserIDs, articleID string) (bool, error) {

	var exists bool
	err := env.dbs.Read.Get(
		&exists,
		pkg.StmtIsStarring,
		userIDs.BuildFindInSet(),
		articleID)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (env StarRepo) countStars(id ids.UserIDs) (int64, error) {
	var count int64
	err := env.dbs.Read.Get(
		&count,
		pkg.StmtCountStarred,
		id.BuildFindInSet(),
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env StarRepo) listStarred(id ids.UserIDs, p gorest.Pagination) ([]pkg.RawStarred, error) {
	var starred = make([]pkg.RawStarred, 0)
	err := env.dbs.Read.Select(
		&starred,
		pkg.StmtListStarred,
		id.BuildFindInSet(),
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		return nil, err
	}

	return starred, nil
}

// ListStarred lists a user's starred articles.
func (env StarRepo) ListStarred(id ids.UserIDs, p gorest.Pagination) (pkg.PagedList[pkg.Starred], error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	countCh := make(chan int64)
	listCh := make(chan pkg.AsyncResult[[]pkg.RawStarred])

	go func() {
		defer close(countCh)
		n, err := env.countStars(id)
		if err != nil {
			sugar.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		l, err := env.listStarred(id, p)
		if err != nil {
			sugar.Error(err)
		}
		listCh <- pkg.AsyncResult[[]pkg.RawStarred]{
			Value: l,
			Err:   err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return pkg.PagedList[pkg.Starred]{}, listResult.Err
	}

	starred := collection.Map(listResult.Value, func(item pkg.RawStarred) pkg.Starred {
		return item.Build()
	})

	return pkg.PagedList[pkg.Starred]{
		Total:      count,
		Pagination: p,
		Data:       starred,
	}, nil
}
