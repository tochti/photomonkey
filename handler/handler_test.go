package handler

import (
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/tochti/photomonkey/database"
	"github.com/tochti/photomonkey/observer"
)

func Test_NextPhoto(t *testing.T) {
	tc := struct {
		Photo database.Photo
	}{
		Photo: database.Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "caption",
		},
	}

	// Run test
	{
		db := database.InitTestDB(t)

		observer := &observer.PhotoObservers{}

		logger := log.New(os.Stdout, "", log.LstdFlags)
		router := NewRouter(db, logger, observer)
		ts := httptest.NewServer(router)

		observer.Broadcast(tc.Photo)

		u := url.URL{
			Scheme: "ws",
			Host:   ts.Listener.Addr().String(),
			Path:   "/v1/new_photos",
		}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()

		photo := database.Photo{}
		err = c.ReadJSON(&photo)
		if err != nil {
			t.Fatal(err)
		}

		if photo.ID != tc.Photo.ID ||
			photo.Hash != tc.Photo.Hash ||
			photo.Caption != tc.Photo.Caption {
			t.Fatalf("Expect %v was %v", tc.Photo, photo)
		}
	}
}
