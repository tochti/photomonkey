package database

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/mattes/migrate/migrate"
)

var (
	TestAppName          = "TEST"
	TestDatabaseConnPool *sql.DB
)

type (
	SQLiteTestConnPool struct {
		Pool           *sql.DB
		MigrationSpecs MigrationSpecs
		SQLiteSpecs    SQLiteSpecs
	}

	TestDatabaseMethods interface {
		IsInTable(table, where string, args ...interface{}) error
		Reset(*testing.T)
	}
)

func (p *SQLiteTestConnPool) IsInTable(table, where string, args ...interface{}) error {

	q := fmt.Sprintf("SELECT * FROM %v WHERE %v", table, where)
	_, err := p.Pool.Query(q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (p *SQLiteTestConnPool) Reset(t *testing.T) {

	errs, ok := migrate.ResetSync(p.SQLiteSpecs.String(), p.MigrationSpecs.Path)
	if !ok {
		t.Fatal(errs)
	}
}

func InitNewTestDB(t *testing.T) (DatabaseMethods, TestDatabaseMethods) {
	migrationSpecs, err := ReadMigrationSpecs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	sqliteSpecs, err := ReadSQLiteSpecs(TestAppName)
	if err != nil {
		t.Fatal(err)
	}

	db := &SQLiteConnPool{TestDatabaseConnPool}
	tDb := &SQLiteTestConnPool{
		Pool:           TestDatabaseConnPool,
		MigrationSpecs: migrationSpecs,
		SQLiteSpecs:    sqliteSpecs,
	}

	tDb.Reset(t)

	return db, tDb
}

func InitSQLiteConnPool(appName string) {
	sqliteSpecs, err := ReadSQLiteSpecs(appName)
	if err != nil {
		log.Fatal(err)
	}

	TestDatabaseConnPool, err = sqliteSpecs.DB()
	if err != nil {
		log.Fatal(err)
	}
}
