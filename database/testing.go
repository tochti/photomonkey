package database

import (
	"crypto/sha1"
	"fmt"
	"testing"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/tochti/dbtt"
	"github.com/tochti/photomonkey/app"
	"github.com/tochti/speci"
)

var (
	TestAppName = app.Name

	TestPhoto = Photo{
		ID:      "AwADBAADYwADO1wlBuF1ogMa7HnMAg",
		Hash:    fmt.Sprintf("%x", sha1.Sum([]byte("42"))),
		Caption: "caption",
	}
)

func InitTestDB(t *testing.T) *PostgreSQLConn {
	migrationSpecs, err := ReadMigrationSpecs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	dbSpecs, err := speci.ReadPostgreSQL(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	db, err := Init("postgres", dbSpecs.String())
	if err != nil {
		t.Fatal(err)
	}

	dbtt.ResetDB(t, dbSpecs.String(), migrationSpecs.Path)

	return db
}
