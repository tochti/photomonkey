package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattes/migrate/driver/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

var (
	TablePhotos = "photos"
)

type (
	DatabaseMethods interface {
		NewPhoto(id, hash, comment string) (Photo, error)
		ReadAllPhotos() ([]Photo, error)
		ReadAllPhotosNewer(time.Time) ([]Photo, error)
		RemovePhotoByHash(hash string) error
	}

	SQLiteConnPool struct {
		Pool *sql.DB
	}

	SQLiteSpecs struct {
		Path string `envconfig:"SQLITE_PATH"`
	}

	MigrationSpecs struct {
		Path string `required:"true" envconfig:"MIGRATION_PATH"`
	}

	Photo struct {
		ID         string
		Hash       string
		Caption    string
		CreateTime time.Time
	}
)

// Insert new photo in database
func (p *SQLiteConnPool) NewPhoto(id, hash, caption string) (Photo, error) {
	q := fmt.Sprintf("INSERT INTO %v VALUES (?,?,?,?)", TablePhotos)

	date := time.Now()
	_, err := p.Pool.Exec(q, id, hash, caption, date)
	if err != nil {
		return Photo{}, err
	}

	photo := Photo{
		ID:         id,
		Hash:       hash,
		Caption:    caption,
		CreateTime: date,
	}
	return photo, nil
}

// Read all photos from database
func (p *SQLiteConnPool) ReadAllPhotos() ([]Photo, error) {
	q := fmt.Sprintf("SELECT * FROM %v", TablePhotos)

	rows, err := p.Pool.Query(q)
	if err != nil {
		return []Photo{}, err
	}

	photos := []Photo{}
	for rows.Next() {
		p := Photo{}
		err := rows.Scan(&p.ID, &p.Hash, &p.Caption, &p.CreateTime)
		if err != nil {
			return []Photo{}, err
		}

		photos = append(photos, p)
	}

	return photos, nil
}

// Read all photos from database newer as X
func (p *SQLiteConnPool) ReadAllPhotosNewer(newer time.Time) ([]Photo, error) {
	q := fmt.Sprintf("SELECT * FROM %v WHERE create_time >= ? ORDER BY create_time", TablePhotos)

	rows, err := p.Pool.Query(q, newer)
	if err != nil {
		return []Photo{}, err
	}

	photos := []Photo{}
	for rows.Next() {
		p := Photo{}
		err := rows.Scan(&p.ID, &p.Hash, &p.Caption, &p.CreateTime)
		if err != nil {
			return []Photo{}, err
		}

		photos = append(photos, p)
	}

	return photos, nil
}

// Remove photo from database by photo hash
func (p *SQLiteConnPool) RemovePhotoByHash(hash string) error {
	q := fmt.Sprintf("DELETE FROM %v WHERE hash=?", TablePhotos)

	_, err := p.Pool.Exec(q, hash)
	if err != nil {
		return err
	}

	return nil
}

func ReadSQLiteSpecs(prefix string) (SQLiteSpecs, error) {
	specs := SQLiteSpecs{}

	err := envconfig.Process(prefix, &specs)
	if err != nil {
		return SQLiteSpecs{}, err
	}

	return specs, nil
}

func (s SQLiteSpecs) DB() (*sql.DB, error) {
	pool, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		return &sql.DB{}, err
	}

	return pool, nil
}

func (s SQLiteSpecs) String() string {
	return fmt.Sprintf("sqlite3://%v", s.Path)
}

func ReadMigrationSpecs(prefix string) (MigrationSpecs, error) {
	specs := MigrationSpecs{}

	err := envconfig.Process(prefix, &specs)
	if err != nil {
		return MigrationSpecs{}, err
	}

	return specs, err
}
