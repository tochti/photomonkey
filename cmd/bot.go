package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/tochti/photomonkey/bot"
	"github.com/tochti/photomonkey/database"
	"github.com/tochti/photomonkey/handler"
	"github.com/tochti/photomonkey/observer"
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

	sqliteSpecs, err := database.ReadSQLiteSpecs(AppName)
	if err != nil {
		logger.Fatal(err)
	}

	pool, err := sqliteSpecs.DB()
	if err != nil {
		logger.Fatal(err)
	}

	db := &database.SQLiteConnPool{pool}

	observers := &observer.PhotoObservers{}

	go bot.Start(logger, observers, db, botSpecs.Token, botSpecs.ImageDir)

	handlers := &handler.Handlers{
		Log:      logger,
		Database: db,
	}

	upgrader := websocket.Upgrader{}

	dir := http.Dir(botSpecs.ImageDir)
	http.Handle("/files", http.StripPrefix("/files/", http.FileServer(dir)))
	//http.HandleFunc("/all_photos", handlers.ReadAllPhotos)
	http.Handle("/new_photos", handlers.ReceiveNewPhotos(upgrader))

	http.ListenAndServe(":8080", nil)
}
