package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FTChinese/go-rest/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/ftchinese/content-api/internal/access"
	"gitlab.com/ftchinese/content-api/internal/controller"
	"gitlab.com/ftchinese/content-api/pkg/config"
	"gitlab.com/ftchinese/content-api/pkg/db"
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
	myDB := db.MustNewMySQL(config.MustMySQLReadConn())

	//accessGuard := controller.AccessGuard{
	//	env: repository.NewOAuthEnv(db),
	//}
	guard := access.NewGuard(myDB)
	storyRoutes := controller.NewStoryRouter(myDB, logger)
	videoRouter := controller.NewVideoRouter(myDB, logger)
	galleryRouter := controller.NewGalleryStory(myDB, logger)
	pageRouter := controller.NewPageRouter(myDB, logger)
	interactiveRouter := controller.NewAudioRouter(myDB, logger)

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
		var data = map[string]string{
			"home":                 "/front-page/latest",
			"home_archive":         "/front-page/archives/{date}",
			"channels_index":       "/channels",
			"channel_page":         "/channels/{name}",
			"story":                "/contents/stories/{id}/<cn | en | ce>",
			"video":                "/contents/videos/{id}",
			"gallery":              "/contents/galleries/{id}",
			"interactive_channels": "/interactive/channels/{name}",
			"interactive_contents": "/interactive/contents/{id}",
		}

		_ = render.New(writer).OK(data)
	})

	r.Route("/front-page", func(r chi.Router) {
		r.Get("/latest", pageRouter.TodayFrontPage)
		r.Get("/archives/{date}", pageRouter.ArchivedFrontPage)
	})

	// GET /channels/{pathName}?page=<int>&per_page=<int>
	r.Route("/channels", func(r chi.Router) {
		// A list of all channel names
		r.Get("/", pageRouter.ChannelList)
		// The details of each channel.
		r.Get("/{name}", pageRouter.ChannelData)
	})

	r.Route("/interactive", func(r chi.Router) {
		// The details of a channel
		r.Get("/channels/{name}", interactiveRouter.ChannelPage)
		// The content of an article.
		r.Get("/contents/{id}", interactiveRouter.Content)
	})

	r.Route("/contents", func(r chi.Router) {
		r.Get("/stories/{id}", storyRoutes.Story)
		r.Get("/videos/{id}", videoRouter.Article)
		r.Get("/galleries/{id}", galleryRouter.Article)
	})

	r.Get("/__version", func(w http.ResponseWriter, req *http.Request) {
		_ = render.New(w).OK(status)
	})

	log.Printf("content-api started on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
