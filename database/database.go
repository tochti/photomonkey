package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
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

	SQLiteConn struct {
		Pool *sqlx.DB
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

// Create new Db connection pool
func Init(driver, url string) (*SQLiteConn, error) {
	pool, err := sqlx.Connect("sqlite3", url)
	if err != nil {
		return &SQLiteConn{}, err
	}

	return &SQLiteConn{
		Pool: pool,
	}, nil
}

// Insert new photo in database
func (p *SQLiteConn) NewPhoto(id, hash, caption string) (Photo, error) {
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
func (p *SQLiteConn) ReadAllPhotos() ([]Photo, error) {
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

// Remove photo from database by photo hash
func (p *SQLiteConn) RemovePhotoByHash(hash string) error {
	q := fmt.Sprintf("DELETE FROM %v WHERE hash=?", TablePhotos)

	_, err := p.Pool.Exec(q, hash)
	if err != nil {
		return err
	}

	return nil
}

func ReadMigrationSpecs(prefix string) (MigrationSpecs, error) {
	specs := MigrationSpecs{}

	err := envconfig.Process(prefix, &specs)
	if err != nil {
		return MigrationSpecs{}, err
	}

	return specs, err
}
