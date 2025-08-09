package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func ApplyMigrationsPG(pgdb *sql.DB) error {
	if err := goose.Up(pgdb, "./migrations/postgeres"); err != nil {
		fmt.Println("Couldn't apply postgres migration: ", err)
		return err
	}
	fmt.Println("Postgress migrations applied successfully!")

	return nil
}
