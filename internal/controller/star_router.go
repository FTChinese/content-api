package controller

import (
	"net/http"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/pkg"
	"github.com/FTchinese/content-api/internal/repository"
	"github.com/FTchinese/content-api/pkg/db"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"go.uber.org/zap"
)

type StarRouter struct {
	repo   repository.StarRepo
	logger *zap.Logger
}

func NewStarRouter(dbs db.ReadWriteMyDBs, logger *zap.Logger) StarRouter {
	return StarRouter{
		repo:   repository.NewStarRepo(dbs, logger),
		logger: logger,
	}
}

// StarredArticles list user's starred articles
//
//	GET /starred?page=<integer>&per_page=<integer>
func (router StarRouter) ListStarred(w http.ResponseWriter, req *http.Request) {
	p := gorest.GetPagination(req)
	userIDs := xhttp.UserIDsFromHeader(req.Header)

	list, err := router.repo.ListStarred(userIDs, p)

	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	_ = render.New(w).OK(list)
}

// IsStarring checks if a user is starring an article
//
//		GET /user/starred/{articleID}
//
//	  - `400 Bad Request` if request URL does not contain `aritcleID` part:
//	    {
//	    "message": "Invalid request URI"
//	    }
//
// * `404 Not Found` is user is not starring this article
//
// * `204 No Content` if user is starring this article
func (router StarRouter) IsStarring(w http.ResponseWriter, req *http.Request) {
	userIDs := xhttp.UserIDsFromHeader(req.Header)

	// Get article id
	articleID, err := xhttp.GetURLParam(req, "storyID").ToString()
	if err != nil {
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	ok, err := router.repo.IsStarring(userIDs, articleID)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	if !ok {
		// 404 Not Found if user is not starring this article
		_ = render.New(w).NotFound("")
		return
	}

	// 204 No Content if user is starring this article
	_ = render.New(w).NoContent()
}

// StarArticle saves an article starred by a user.
//
//		PUT /user/starred/{articleID}
//
//	  - `400 Bad Request` if request URL does not contain `aritcleID` part:
//	    {
//	    "message": "Invalid request URI"
//	    }
//
// * Response `204 No Content`
//
// This endpoint always returns `204 No Content`. Duplicate record will be simply discarded.
func (router StarRouter) StarArticle(w http.ResponseWriter, req *http.Request) {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	userIDs := xhttp.UserIDsFromHeader(req.Header)

	var params pkg.StarParams
	if err := gorest.ParseJSON(req.Body, &params); err != nil {
		sugar.Error(err)
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	if ve := params.Validate(); ve != nil {
		sugar.Error(ve)
		_ = render.New(w).Unprocessable(ve)
	}

	params.UserID = userIDs.GetCompoundID()
	err := router.repo.StarArticle(params)
	if err != nil {
		_ = render.New(w).DBError(err)
		return
	}

	// Respond 204
	_ = render.New(w).NoContent()
}

// UnstarArticle delete a piece of user starred article.
//
//		DELETE /user/starred/{articleID}
//
//	  - `400 Bad Request` if request URL does not contain `aritcleID` part:
//	    {
//	    "message": "Invalid request URI"
//	    }
//
// * Response `204 No Content`.
func (router StarRouter) UnstarArticle(w http.ResponseWriter, req *http.Request) {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	userIDs := xhttp.UserIDsFromHeader(req.Header)

	// Get article id
	articleID, err := xhttp.GetURLParam(req, "storyID").ToString()
	if err != nil {
		sugar.Error(err)
		_ = render.New(w).BadRequest(err.Error())
		return
	}

	// Delete
	err = router.repo.UnstarArticle(userIDs.GetCompoundID(), articleID)
	if err != nil {
		sugar.Error(err)
		_ = render.New(w).DBError(err)
		return
	}

	// Respond 204
	_ = render.New(w).NoContent()
}
