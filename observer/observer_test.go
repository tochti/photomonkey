package observer

import (
	"testing"

	"github.com/tochti/photomonkey/database"
)

func Test_Observer(t *testing.T) {
	observer := New()

	client_channel1 := make(chan database.Photo)
	client_channel2 := make(chan database.Photo)

	id1 := observer.Add(client_channel1)
	id2 := observer.Add(client_channel2)

	tc := struct {
		Expected database.Photo
	}{
		Expected: database.Photo{
			ID:      "123",
			Hash:    "love",
			Caption: "and peace",
		},
	}
	observer.Broadcast(tc.Expected)

	p := <-client_channel1
	if tc.Expected.ID != p.ID ||
		tc.Expected.Hash != p.Hash ||
		tc.Expected.Caption != p.Caption {
		t.Fatalf("Expect %v was %v", tc.Expected, p)
	}

	p = <-client_channel2
	if tc.Expected.ID != p.ID ||
		tc.Expected.Hash != p.Hash ||
		tc.Expected.Caption != p.Caption {
		t.Fatalf("Expect %v was %v", tc.Expected, p)
	}

	observer.Remove(id1)
	observer.Remove(id2)
}
