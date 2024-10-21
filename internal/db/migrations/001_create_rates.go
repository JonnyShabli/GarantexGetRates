package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(UP_001, Down_001)
}

func UP_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    ts int NOT NULL,
    ask_price VARCHAR(15) NOT NULL,
    ask_volume VARCHAR(15) NOT NULL,
    ask_amount VARCHAR(15) NOT NULL,
    ask_factor VARCHAR(15) NOT NULL,
    ask_type VARCHAR(15) NOT NULL,
    bid_price VARCHAR(15) NOT NULL,
    bid_volume VARCHAR(15) NOT NULL,
    bid_amount VARCHAR(15) NOT NULL,
    bid_factor VARCHAR(15) NOT NULL,
    bid_type VARCHAR(15) NOT NULL                             
);`)
	if err != nil {
		return err
	}
	return nil
}

func Down_001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE rates;`)
	if err != nil {
		return err
	}
	return nil
}
