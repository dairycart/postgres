package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	dairycart "github.com/dairycart/dairycart/api"
	"github.com/dairycart/postgres/migrations"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
)

const (
	maxConnectionAttempts = 5
)

func loadMigrationData(dbURL string, loadExampleData bool) (*migrate.Migrate, error) {
	s := bindata.Resource(migrations.AssetNames(), func(name string) ([]byte, error) {
		if strings.Contains(name, "example_data") && loadExampleData {
			return migrations.Asset(name)
		}
		return migrations.Asset(name)
	})
	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithSourceInstance("go-bindata", d, dbURL)
}

func prepareForMigration(db *sql.DB, dbURL string, loadExampleData bool) (*migrate.Migrate, error) {
	log.Printf("preparing to migrate postgres database at url: '%s'\n", dbURL)
	err := databaseIsAvailable(db)
	if err != nil {
		return nil, err
	}

	m, err := loadMigrationData(dbURL, loadExampleData)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func databaseIsAvailable(db *sql.DB) error {
	numberOfUnsuccessfulAttempts := 0
	databaseIsNotMigrated := true
	for databaseIsNotMigrated {
		err := db.Ping()
		if err != nil {
			log.Printf("ping failed, waiting half a second for the database")
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

func (pg *postgres) Migrate(db *sql.DB, cfg *viper.Viper) error {
	dbURL := cfg.GetString(dairycart.DatabaseConnectionKey)
	loadExampleData := cfg.GetBool(dairycart.MigrateExampleDataKey)

	m, err := prepareForMigration(db, dbURL, false)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

func (pg *postgres) Downgrade(db *sql.DB, dbURL string) error {
	m, err := prepareForMigration(db, dbURL, false)
	if err != nil {
		return err
	}

	err = m.Down()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
