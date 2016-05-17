package handler

import (
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tochti/photomonkey/database"
)

func Test_NextPhoto(t *testing.T) {
	tc := struct {
		Photos []database.Photo
	}{
		Photos: []database.Photo{
			{
				ID:      "123",
				Hash:    "hash",
				Caption: "caption",
			},
		},
	}

	db, _ := database.InitNewTestDB(t)

	upgrader := websocket.Upgrader{}

	null := os.NewFile(uintptr(syscall.Stdin), os.DevNull)
	handlers := &Handlers{
		Log:      log.New(null, "", log.LstdFlags),
		Database: db,
	}

	newPhoto := handlers.ReceiveNewPhotos(upgrader)
	ts := httptest.NewServer(newPhoto)

	time.AfterFunc(200*time.Millisecond, func() {
		for _, photo := range tc.Photos {
			_, err := db.NewPhoto(photo.ID, photo.Hash, photo.Caption)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	u := url.URL{
		Scheme: "ws",
		Host:   ts.Listener.Addr().String(),
		Path:   "/",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	photos := []database.Photo{}
	if err != nil {
		t.Fatal(err)
	}
	err = c.ReadJSON(&photos)
	if err != nil {
		t.Fatal(err)
	}

	if len(photos) != len(tc.Photos) {
		t.Fatalf("Expect len %v was %v", len(tc.Photos), len(photos))
	}

	for i, photo := range tc.Photos {
		if photos[i].ID != photo.ID ||
			photos[i].Hash != photo.Hash ||
			photos[i].Caption != photo.Caption {
			t.Fatalf("Expect %v was %v", photo, photos[i])
		}
	}

}

func init() {
	database.InitSQLiteConnPool(database.TestAppName)
}
