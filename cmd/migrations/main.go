package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/chrisdiebold/snippetbox/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func runMigrations(databaseURL string) error {
	d, err := iofs.New(migrations.MigrationFS, ".")
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("migrations complete")
	return nil
}

func main() {
	// TODO: make this so it does not leak default values
	user := flag.String("user", "developer", "Database user name")
	password := flag.String("password", "developer", "Database password")

	flag.Parse()

	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/snippetbox_dev?sslmode=disable", *user, *password)

	if err := runMigrations(connStr); err != nil {
		log.Fatal(err)
	}

}
