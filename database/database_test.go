package database

import (
	"database/sql"
	"testing"
	"time"
)

func Test_NewPhoto(t *testing.T) {
	db, tDb := InitNewTestDB(t)

	tc := struct {
		Photo Photo
	}{
		Photo: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "caption",
		},
	}

	photo, err := db.NewPhoto(tc.Photo.ID, tc.Photo.Hash, tc.Photo.Caption)
	if err != nil {
		t.Fatal(err)
	}

	if photo.ID != tc.Photo.ID ||
		photo.Hash != tc.Photo.Hash ||
		photo.Caption != tc.Photo.Caption {
		t.Fatalf("Expect %v was %v", tc.Photo, photo)
	}

	err = tDb.IsInTable(
		TablePhotos,
		"id=? AND hash=?",
		tc.Photo.ID, tc.Photo.Hash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatalf("Expect %v in database", tc.Photo)
		}

		t.Fatal(err)
	}

}

func Test_ReadAllPhotos(t *testing.T) {
	db, _ := InitNewTestDB(t)

	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

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

func Test_ReadAllPhotosNewer(t *testing.T) {
	db, _ := InitNewTestDB(t)

	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	tmp, err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
	if err != nil {
		t.Fatal(err)
	}

	photos, err := db.ReadAllPhotosNewer(tmp.CreateTime.Add(-1 * time.Nanosecond))

	if len(photos) != 1 {
		t.Fatalf("Expect len %v was %v", 1, len(photos))
	}

	if tc.Expected.ID != photos[0].ID ||
		tc.Expected.Hash != photos[0].Hash ||
		tc.Expected.Caption != photos[0].Caption {
		t.Fatalf("Expect %v was %v", tc.Expected, photos)
	}
}

func Test_RemovePhoto(t *testing.T) {
	db, tDb := InitNewTestDB(t)

	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	_, err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
	if err != nil {
		t.Fatal(err)
	}

	err = db.RemovePhotoByHash(tc.Expected.Hash)
	if err != nil {
		t.Fatal(err)
	}

	err = tDb.IsInTable(
		TablePhotos,
		"id=? AND hash=?",
		tc.Expected.ID, tc.Expected.Hash,
	)

	if err != nil && err != sql.ErrNoRows {
		t.Fatal(err)
	}

}

func init() {
	InitSQLiteConnPool(TestAppName)
}
