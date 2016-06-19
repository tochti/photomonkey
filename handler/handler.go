package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/tochti/hrr"
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
		Log      *logrus.Logger
		PhotoC   chan database.Photo
	}
)

func NewRouter(db database.DatabaseMethods, log *logrus.Logger, observers *observer.PhotoObservers) *httprouter.Router {
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
	router.Handler("GET", "/v1/photos", http.HandlerFunc(handler.ReadAllPhotos))

	return router
}

func (ctx *Handlers) ReceiveNewPhotos(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			e := hrr.NewError("Cannot start websocket connection", err)
			hrr.Response(w, r).Error(e)
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
		updateTime = time.Now()
		go func() {
			ws.WriteJSON(photo)
		}()
	}
}

func (ctx *Handlers) ReadAllPhotos(w http.ResponseWriter, r *http.Request) {
	if err := hrr.Request(r).Log().Process(); err != nil {
		hrr.Response(w, r).Error(err)
		return
	}

	hrr.Response(w, r).Data(func() (interface{}, hrr.Error) {
		photos, err := ctx.Database.ReadAllPhotos()
		if err != nil {
			return nil, hrr.NewError("Cannot read all photos", err)
		}

		return photos, nil
	})
}

func ErrorResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := ErrorMessage{msg}
	json.NewEncoder(w).Encode(err)
}
