package handler

import (
	"net/http/httptest"
	"net/url"
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
			{
				ID:      "123",
				Hash:    "hash",
				Caption: "caption",
			},
		},
	}

	// Run test
	{
		db := database.InitTestDB(t)

		upgrader := websocket.Upgrader{}
		observer := &observer.PhotoObservers{}

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

		photos := []database.Photo{}
		err = c.ReadJSON(&photos)
		if err != nil {
			t.Fatal(err)
		}

		if len(photos) != 1 {
			t.Fatalf("Expect len %v was %v", 1, len(photos))
		}

		if photos[0].ID != tc.Photo.ID ||
			photos[0].Hash != tc.Photo.Hash ||
			photos[0].Caption != tc.Photo.Caption {
			t.Fatalf("Expect %v was %v", tc.Photo, photos[0])
		}
	}
}
