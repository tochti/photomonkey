package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
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
		Database  database.DatabaseMethods
		Log       *logrus.Logger
		Observers *observer.PhotoObservers
	}
)

func NewRouter(db database.DatabaseMethods, log *logrus.Logger, observers *observer.PhotoObservers) http.Handler {
	router := http.NewServeMux()

	hrr.Logger = log

	handler := Handlers{
		Database:  db,
		Log:       log,
		Observers: observers,
	}
	upgrader := websocket.Upgrader{}

	router.HandleFunc("/v1/new_photos", handler.ReceiveNewPhotos(upgrader))
	router.HandleFunc("/v1/photos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		handler.ReadAllPhotos(w, r)
	})

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

		ctx.servePhotos(ws)
	}

}

func (ctx *Handlers) servePhotos(ws *websocket.Conn) {
	ctx.Log.Println("Start to serve photos")

	photoC := make(chan database.Photo)
	chanID := ctx.Observers.Add(photoC)

	go func() {
		defer func() {
			ctx.Log.Println("Close connection ", chanID)
			ctx.Observers.Remove(chanID)
			ws.Close()
		}()

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					ctx.Log.Error(err)
				}

				break
			}
		}
	}()

	for photo := range photoC {
		go func() {
			l := fmt.Sprintf("Share new photo with websocket %v", chanID)
			ctx.Log.Println(l)
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
