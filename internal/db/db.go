package db

import (
	"database/sql"
	"log"

	// Запускает миграции.
	_ "github.com/JonnyShabli/GarantexGetRates/internal/db/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	// драйвер для sqlx.
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func NewConn(connString string) *sqlx.DB {
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return conn
}

func InitDB(connString string) error {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("Starting init DB")

	err = goose.Up(conn, ".")

	if err != nil {
		return err
	}

	return nil
}
