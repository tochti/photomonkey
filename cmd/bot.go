package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/tochti/photomonkey/app"
	"github.com/tochti/photomonkey/bot"
	"github.com/tochti/photomonkey/database"
	"github.com/tochti/photomonkey/handler"
	"github.com/tochti/photomonkey/observer"
	"github.com/tochti/speci"
)

const (
	AppName = app.Name
)

type (
	FrontendSpecs struct {
		Dir string `envconfig:"FRONTEND_DIR" required:"true"`
	}
)

func main() {

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	botSpecs, err := bot.ReadSpecs(AppName)
	if err != nil {
		logger.Fatal(err)
	}

	httpServerSpecs, err := speci.ReadHTTPServer(AppName)
	if err != nil {
		logger.Println(err)
	}

	sqlSpecs, err := speci.ReadPostgreSQL(AppName)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := database.Init("postgres", sqlSpecs.String())
	if err != nil {
		logger.Fatal(err)
	}

	observers := observer.New()

	go bot.Start(logger, observers, db, botSpecs.Token, botSpecs.ImageDir)

	router := handler.NewRouter(db, logger, observers)

	initFrontendRouter(logger)

	dir := http.Dir(botSpecs.ImageDir)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(dir)))
	http.Handle("/v1/", router)

	logger.Println("Listen on " + httpServerSpecs.String())
	http.ListenAndServe(httpServerSpecs.String(), nil)
}

func initFrontendRouter(logger *logrus.Logger) {
	specs := FrontendSpecs{}
	err := envconfig.Process(AppName, &specs)
	if err != nil {
		logger.Fatalln(err)
	}

	fs := http.FileServer(http.Dir(specs.Dir))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
}
