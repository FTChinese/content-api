package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTchinese/content-api/internal/access"
	"github.com/FTchinese/content-api/internal/controller"
	"github.com/FTchinese/content-api/pkg/config"
	"github.com/FTchinese/content-api/pkg/db"
	"github.com/FTchinese/content-api/pkg/xhttp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//go:embed build/api.toml
var tomlConfig string

//go:embed build/version
var version string

//go:embed build/build_time
var build string

//go:embed build/commit
var commit string

var (
	production bool
)

const port = "8100"

func init() {
	flag.BoolVar(&production, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper([]byte(tomlConfig))
}

func main() {
	status := config.ServerStatus{
		Version:    version,
		BuiltAt:    build,
		Commit:     commit,
		Production: production,
	}

	logger := config.MustGetLogger(production)
	myDBs := db.MustNewMyDBs()

	//accessGuard := controller.AccessGuard{
	//	env: repository.NewOAuthEnv(db),
	//}
	guard := access.NewGuard(myDBs.Read)
	storyRoutes := controller.NewStoryRouter(myDBs.Read, logger)
	videoRouter := controller.NewVideoRouter(myDBs.Read, logger)
	galleryRouter := controller.NewGalleryStory(myDBs.Read, logger)
	pageRouter := controller.NewPageRouter(myDBs.Read, logger)
	starRouter := controller.NewStarRouter(myDBs, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(controller.LogRequest)

	r.Use(guard.CheckToken)

	r.Route("/__status", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			var data = map[string]string{
				"channelIds": "/__status/channel_ids",
			}

			_ = render.New(writer).OK(data)
		})
		r.Get("/channel_ids", pageRouter.InspectChannelMap)
	})

	// Sitemap
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		r.Use(guard.CheckToken)
		var data = map[string]string{
			"home":                 "/front-page/latest",
			"home_archive":         "/front-page/archives/{date}",
			"channels_index":       "/channels",
			"channel_page":         "/channels/{name}",
			"story":                "/contents/stories/{id}",
			"video":                "/contents/videos/{id}",
			"gallery":              "/contents/galleries/{id}",
			"interactive_channels": "/interactive/channels/{name}",
			"interactive_contents": "/interactive/contents/{id}",
		}

		_ = render.New(writer).OK(data)
	})

	r.Route("/front-page", func(r chi.Router) {
		r.Use(guard.CheckToken)
		r.Get("/latest", pageRouter.TodayFrontPage)
		r.Get("/archives/{date}", pageRouter.ArchivedFrontPage)
	})

	// GET /channels/{pathName}?page=<int>&per_page=<int>
	r.Route("/channels", func(r chi.Router) {
		r.Use(guard.CheckToken)
		// A list of all channel names
		r.Get("/", pageRouter.ChannelList)
		// The details of each channel.
		r.Get("/{name}", pageRouter.ChannelData)
	})

	// r.Route("/interactive", func(r chi.Router) {
	// 	// The details of a channel
	// 	r.Get("/channels/{name}", interactiveRouter.ChannelPage)
	// 	// The content of an article.
	// 	r.Get("/contents/{id}", interactiveRouter.Content)
	// })

	r.Route("/contents", func(r chi.Router) {
		r.Use(guard.CheckToken)
		r.Get("/stories/{id}", storyRoutes.Story)
		r.Get("/videos/{id}", videoRouter.Article)
		r.Get("/galleries/{id}", galleryRouter.Article)
	})

	r.Route("/starred", func(r chi.Router) {
		r.Use(guard.CheckToken)
		r.Use(xhttp.RequireFtcOrUnionID)

		// ?page=<int>&per_page=<int>
		r.With(xhttp.FormParsed).Get("/", starRouter.ListStarred)
		r.Post("/", starRouter.StarArticle)
		r.Get("/{storyID}", starRouter.IsStarring)
		r.Delete("/{storyID}", starRouter.UnstarArticle)
	})

	r.Get("/__version", func(w http.ResponseWriter, req *http.Request) {
		_ = render.New(w).OK(status)
	})

	log.Printf("content-api started on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
