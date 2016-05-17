package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tochti/photomonkey/database"
)

var (
	UpdateInterval = 2 * time.Second
)

type (
	ErrorMessage struct {
		Message string `json:"message"`
	}

	Handlers struct {
		Log      *log.Logger
		Database database.DatabaseMethods
	}
)

func (ctx *Handlers) ReceiveNewPhotos(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ticker := time.NewTicker(time.Second)
		go ctx.servePhotos(ws, ticker)
	}

}

func (ctx *Handlers) servePhotos(ws *websocket.Conn, ticker *time.Ticker) {
	ctx.Log.Println("Start to serve photos")
	updateTime := time.Now()
	for range ticker.C {
		ctx.Log.Println("Last Update Time:", updateTime)
		ctx.Log.Println("Looking for updates")
		photos, err := ctx.Database.ReadAllPhotosNewer(updateTime)
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
