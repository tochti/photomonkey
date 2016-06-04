package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tochti/photomonkey/bot"
	"github.com/tochti/photomonkey/database"
	"github.com/tochti/photomonkey/handler"
	"github.com/tochti/photomonkey/observer"
	"github.com/tochti/speci"
)

const (
	AppName = "Photomonkey"
)

func main() {

	logger := log.New(os.Stdout, fmt.Sprintf("%v: ", AppName), log.LstdFlags)

	botSpecs, err := bot.ReadSpecs()
	if err != nil {
		logger.Fatal(err)
	}

	httpServerSpecs, err := speci.ReadHTTPServer(AppName)
	if err != nil {
		logger.Println(err)
	}

	sqliteSpecs, err := speci.ReadSQLite(AppName)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := database.Init("sqlite3", sqliteSpecs.String())
	if err != nil {
		logger.Fatal(err)
	}

	observers := &observer.PhotoObservers{}

	go bot.Start(logger, observers, db, botSpecs.Token, botSpecs.ImageDir)

	router := handler.NewRouter(db, logger, observers)

	dir := http.Dir(botSpecs.ImageDir)
	http.Handle("/files", http.StripPrefix("/files/", http.FileServer(dir)))
	http.Handle("/v1", router)

	http.ListenAndServe(httpServerSpecs.String(), nil)
}
