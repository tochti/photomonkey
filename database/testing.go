package database

import (
	"testing"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/tochti/dbtt"
	"github.com/tochti/speci"
)

var (
	TestAppName = "TEST"
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
