package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
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

	PostgreSQLConn struct {
		Pool *sqlx.DB
	}

	MigrationSpecs struct {
		Path string `required:"true" envconfig:"MIGRATION_PATH"`
	}

	Photo struct {
		ID         string    `db:"id"`
		Hash       string    `db:"hash"`
		Caption    string    `db:"caption"`
		CreateTime time.Time `db:"create_time"`
	}
)

// Create new Db connection pool
func Init(driver, url string) (*PostgreSQLConn, error) {
	pool, err := sqlx.Connect(driver, url)
	if err != nil {
		return &PostgreSQLConn{}, err
	}

	return &PostgreSQLConn{
		Pool: pool,
	}, nil
}

// Insert new photo in database
func (p *PostgreSQLConn) NewPhoto(id, hash, caption string) (Photo, error) {
	q := p.Pool.Rebind(fmt.Sprintf("INSERT INTO %v VALUES (?,?,?) RETURNING *", TablePhotos))

	photo := Photo{}
	err := p.Pool.Get(&photo, q, id, hash, caption)
	if err != nil {
		return Photo{}, err
	}

	return photo, nil
}

// Read all photos from database
func (p *PostgreSQLConn) ReadAllPhotos() ([]Photo, error) {
	q := p.Pool.Rebind(fmt.Sprintf("SELECT * FROM %v", TablePhotos))

	photos := []Photo{}
	err := p.Pool.Select(&photos, q)
	if err != nil {
		return []Photo{}, err
	}

	return photos, nil
}

// Remove photo from database by photo hash
func (p *PostgreSQLConn) RemovePhotoByHash(hash string) error {
	q := p.Pool.Rebind(fmt.Sprintf("DELETE FROM %v WHERE hash=?", TablePhotos))

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
