package database

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/mattes/migrate/migrate"
)

var (
	TestDatabaseConnPool *sql.DB
	TestAppName          = "TEST"
)

type (
	TestDatabaseMethods interface {
		IsInTable(table, where string, args ...interface{}) error
		Reset(*testing.T)
	}

	SQLiteTestConnPool struct {
		Pool *sql.DB
	}
)

func Test_NewPhoto(t *testing.T) {
	db, tDb := initNewTestDB(t)

	tc := struct {
		Photo Photo
	}{
		Photo: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "caption",
		},
	}

	err := db.NewPhoto(tc.Photo.ID, tc.Photo.Hash, tc.Photo.Caption)
	if err != nil {
		t.Fatal(err)
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
	db, _ := initNewTestDB(t)

	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
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

func Test_RemovePhoto(t *testing.T) {
	db, tDb := initNewTestDB(t)

	tc := struct {
		Expected Photo
	}{
		Expected: Photo{
			ID:      "123",
			Hash:    "hash",
			Caption: "comment",
		},
	}

	err := db.NewPhoto(tc.Expected.ID, tc.Expected.Hash, tc.Expected.Caption)
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

func (p *SQLiteTestConnPool) IsInTable(table, where string, args ...interface{}) error {

	q := fmt.Sprintf("SELECT * FROM %v WHERE %v", table, where)
	_, err := p.Pool.Query(q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (*SQLiteTestConnPool) Reset(t *testing.T) {
	migrationSpecs, err := ReadMigrationSepcs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	sqliteSpecs, err := ReadSQLiteSpecs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	errs, ok := migrate.ResetSync(sqliteSpecs.String(), migrationSpecs.Path)
	if !ok {
		t.Fatal(errs)
	}
}

func initNewTestDB(t *testing.T) (DatabaseMethods, TestDatabaseMethods) {
	db := &SQLiteConnPool{TestDatabaseConnPool}
	tDb := &SQLiteTestConnPool{TestDatabaseConnPool}

	tDb.Reset(t)

	return db, tDb
}

func init() {
	sqliteSpecs, err := ReadSQLiteSpecs(TestAppName)
	if err != nil {
		log.Fatal(err)
	}

	TestDatabaseConnPool, err = sqliteSpecs.DB()
	if err != nil {
		log.Fatal(err)
	}
}
