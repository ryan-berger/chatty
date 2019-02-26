package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

const MIGRATIONS_PATH = "file:///app/repositories/postgres/migrations"

func GetDBString() string {
	return fmt.Sprintf(
		"host=%s database=%s user=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"))
}

func main() {
	db, err := sqlx.Open("postgres", GetDBString())

	if err != nil {
		panic(err)
	}

	fmt.Println("migrating")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("driver", driver)

	m, err := migrate.NewWithDatabaseInstance(
		MIGRATIONS_PATH,
		"postgres",
		driver,
	)

	if err != nil {
		panic(err)
	}

	err = m.Up()

	if err != nil {
		panic(err)
	}
}
