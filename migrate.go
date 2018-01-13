package postgres

import (
	"strings"

	"github.com/dairycart/postgres/migrations"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
)

const (
	maxConnectionAttempts = 25
)

func loadMigrations(dbURL string) (*migrate.Migrate, error) {
	s := bindata.Resource(migrations.AssetNames(), func(name string) ([]byte, error) {
		if !strings.Contains(name, "example_data") {
			return migrations.Asset(name)
		}
		return nil, nil
	})
	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("go-bindata", d, dbURL)
}

func (pg *postgres) Migrate(dbURL string) error {
	m, err := loadMigrations(dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

func (pg *postgres) Downgrade(dbURL string) error {
	m, err := loadMigrations(dbURL)
	if err != nil {
		return err
	}

	err = m.Down()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
