package repository

import (
	"context"
	"fmt"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type GgrRepoInterface interface {
	InsertRates(ctx context.Context, db models.RatesToDB) error
}

type GgrRepoObj struct {
	log *zap.Logger
	db  *sqlx.DB
}

func NewGgrRepo(log *zap.Logger, conn *sqlx.DB) GgrRepoInterface {
	return &GgrRepoObj{
		log: log,
		db:  conn,
	}
}

func (g *GgrRepoObj) InsertRates(ctx context.Context, data models.RatesToDB) error {
	sqlstring, args, err := squirrel.Insert("rates").
		Columns("ts",
			"ask_price",
			"ask_volume",
			"ask_amount",
			"ask_factor",
			"ask_type",
			"bid_price",
			"bid_volume",
			"bid_amount",
			"bid_factor",
			"bid_type",
		).
		Values(data.Timestamp,
			data.Ask.Price,
			data.Ask.Volume,
			data.Ask.Amount,
			data.Ask.Factor,
			data.Ask.Type,
			data.Bid.Price,
			data.Bid.Volume,
			data.Bid.Amount,
			data.Bid.Factor,
			data.Bid.Type).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		g.log.Error("InsertRates", zap.Error(err))
		return fmt.Errorf("InsertRates: %w", err)
	}

	g.log.Info("InsertRates", zap.String("sql", sqlstring))

	_, err = g.db.ExecContext(ctx, sqlstring, args...)
	if err != nil {
		g.log.Error("InsertRates", zap.Error(err))
		return fmt.Errorf("InsertRates: %w", err)
	}
	return nil
}
