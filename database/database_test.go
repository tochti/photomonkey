package database

import (
	"testing"

	"github.com/tochti/dbtt"
)

func Test_NewPhoto(t *testing.T) {
	tc := struct {
		Photo Photo
	}{
		Photo: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "caption",
		},
	}

	// Run test
	{
		db := InitTestDB(t)
		defer db.Pool.Close()

		photo, err := db.NewPhoto(tc.Photo.ID, tc.Photo.Hash, tc.Photo.Caption)
		if err != nil {
			t.Fatal(err)
		}

		if photo.ID != tc.Photo.ID ||
			photo.Hash != tc.Photo.Hash ||
			photo.Caption != tc.Photo.Caption {
			t.Fatalf("Expect %v was %v", tc.Photo, photo)
		}

		dbtt.IsInTable(
			t,
			db.Pool,
			TablePhotos,
			"id=? AND hash=?",
			tc.Photo.ID, tc.Photo.Hash,
		)

	}
}

func Test_ReadAllPhotos(t *testing.T) {
	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	// Run test
	{
		db := InitTestDB(t)
		_, err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
		if err != nil {
			t.Fatal(err)
		}

		photos, err := db.ReadAllPhotos()

		if len(photos) != 1 {
			t.Fatalf("Expect len %v was %v", 1, len(photos))
		}

		if tc.Expected.ID != photos[0].ID ||
			tc.Expected.Hash != photos[0].Hash ||
			tc.Expected.Caption != photos[0].Caption {
			t.Fatalf("Expect %v was %v", tc.Expected, photos)
		}
	}
}

func Test_RemovePhoto(t *testing.T) {
	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	// Run test
	{
		db := InitTestDB(t)
		_, err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
		if err != nil {
			t.Fatal(err)
		}

		err = db.RemovePhotoByHash(tc.Expected.Hash)
		if err != nil {
			t.Fatal(err)
		}

		dbtt.IsNotInTable(
			t,
			db.Pool,
			TablePhotos,
			"id=? AND hash=?",
			tc.Expected.ID, tc.Expected.Hash,
		)

	}

}
