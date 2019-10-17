package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/go-rest/view"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/content-api/config"
	"gitlab.com/ftchinese/content-api/controller"
	"gitlab.com/ftchinese/content-api/repository"
	"net/http"
	"os"
)

var (
	buildConfig config.BuildConfig
	version     string
	buildAt     string
)

const port = "8100"

func init() {
	flag.BoolVar(&buildConfig.Production, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, buildAt)
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	if err := viper.ReadInConfig(); err != nil {
		os.Exit(1)
	}
}

func main() {
	db, err := repository.NewDB(getDBConn())
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	accessGuard := controller.AccessGuard{
		Env: repository.NewOAuthEnv(db),
	}
	storyRouter := controller.NewStoryRouter(db)
	videoRouter := controller.NewVideoRouter(db)
	galleryRouter := controller.NewGalleryStory(db)
	channelRouter := controller.NewChannelRouter(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(controller.LogRequest)

	r.Use(accessGuard.CheckToken)

	r.Get("/__version", func(writer http.ResponseWriter, request *http.Request) {
		_ = view.Render(writer, view.NewResponse().SetBody(buildConfig))
	})

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		var data = map[string]string{
			"story":   "/story/{id}/<cn | en | ce>",
			"video":   "/video/{id}",
			"gallery": "/gallery/{id}",
		}

		_ = view.Render(writer, view.NewResponse().SetBody(data))
	})

	r.Route("/front_page", func(r chi.Router) {
		r.Get("/latest", channelRouter.TodayFrontPage)
		r.Get("/archive/{date}", channelRouter.ArchivedFrontPage)
	})

	r.Route("/story/{id}", func(r chi.Router) {
		r.Get("/", storyRouter.Raw)
		r.Get("/cn", storyRouter.CN)
		r.Get("/en", storyRouter.EN)
		r.Get("/ce", storyRouter.Bilingual)
	})

	r.Get("/video/{id}", videoRouter.Article)

	r.Get("/gallery/{id}", galleryRouter.Article)

	logrus.Infof("content-api started on port %s", port)

	logrus.Fatal(http.ListenAndServe(":"+port, r))
}

func getDBConn() repository.Conn {
	var conn repository.Conn
	var key = "mysql.dev"
	if buildConfig.Production {
		key = "mysql.master"
	}

	if err := viper.UnmarshalKey(key, &conn); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	logrus.Infof("Using MySQL server %s", conn.Host)

	return conn
}
