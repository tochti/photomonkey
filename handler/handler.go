package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/tochti/photomonkey/database"
	"github.com/tochti/photomonkey/observer"
)

var (
	UpdateInterval = 2 * time.Second
)

type (
	ErrorMessage struct {
		Message string `json:"message"`
	}

	Handlers struct {
		Database database.DatabaseMethods
		Log      *log.Logger
		PhotoC   chan database.Photo
	}
)

func NewRouter(db database.DatabaseMethods, log *log.Logger, observers *observer.PhotoObservers) *httprouter.Router {
	router := httprouter.New()

	photoC := make(chan database.Photo)
	observers.Add(photoC)

	handler := Handlers{
		Database: db,
		Log:      log,
		PhotoC:   photoC,
	}
	upgrader := websocket.Upgrader{}

	router.Handler("GET", "/v1/new_photos", handler.ReceiveNewPhotos(upgrader))
	router.Hanlder("GET", "/v1/photos", handler.ReadAllPhotos())

	return router
}

func (ctx *Handlers) ReceiveNewPhotos(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		go ctx.servePhotos(ws)
	}

}

func (ctx *Handlers) servePhotos(ws *websocket.Conn) {
	ctx.Log.Println("Start to serve photos")
	updateTime := time.Now()
	for photo := range ctx.PhotoC {
		ctx.Log.Println("Last Update Time:", updateTime)
		ctx.Log.Println("Looking for updates")
		if err != nil {
			ctx.Log.Println("Error:", err)
			continue
		}
		updateTime = time.Now()
		ws.WriteJSON(photos)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := ErrorMessage{msg}
	json.NewEncoder(w).Encode(err)
}
