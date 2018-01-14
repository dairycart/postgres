package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dairycart/postgres/migrations"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
)

const (
	maxConnectionAttempts = 5
)

func loadMigrationData(dbURL string) (*migrate.Migrate, error) {
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

func prepareForMigration(dbURL string) (*migrate.Migrate, error) {
	m, err := loadMigrationData(dbURL)
	if err != nil {
		return nil, err
	}

	err = databaseIsAvailable(dbURL)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func databaseIsAvailable(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	numberOfUnsuccessfulAttempts := 0
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		err = db.Ping()
		if err != nil {
			log.Printf("waiting half a second for the database")
			time.Sleep(500 * time.Millisecond)
			numberOfUnsuccessfulAttempts++

			if numberOfUnsuccessfulAttempts == maxConnectionAttempts {
				return fmt.Errorf("failed to connect to the database: %v\n", err)
			}
		} else {
			break
		}
	}
	return nil
}

func (pg *postgres) Migrate(dbURL string) error {
	m, err := prepareForMigration(dbURL)
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
	m, err := prepareForMigration(dbURL)
	if err != nil {
		return err
	}

	err = m.Down()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
