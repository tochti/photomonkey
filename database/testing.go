package database

import (
	"testing"

	"github.com/tochti/dbtt"
	"github.com/tochti/speci"

	_ "github.com/mattes/migrate/driver/sqlite3"
)

var (
	TestAppName = "TEST"
)

func InitTestDB(t *testing.T) *SQLiteConn {
	migrationSpecs, err := ReadMigrationSpecs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	sqliteSpecs, err := speci.ReadSQLite(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	db, err := Init("sqlite3", sqliteSpecs.String())
	if err != nil {
		t.Fatal(err)
	}

	dbtt.ResetDB(t, sqliteSpecs.String(), migrationSpecs.Path)

	return db
}
