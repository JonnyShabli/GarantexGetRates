package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"go.uber.org/zap"
)

type GgrServiceInterface interface {
	GetRates(ctx context.Context, pair string) (models.GarantexRates, error)
}

type GgrServiceObj struct {
	log *zap.Logger
}

func NewGgrService(log *zap.Logger) GgrServiceInterface {
	return &GgrServiceObj{log: log}
}

func (g *GgrServiceObj) GetRates(ctx context.Context, pair string) (models.GarantexRates, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://garantex.org/api/v2/depth?market=%s", pair)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		g.log.Info("error forming request")
		return models.GarantexRates{}, err
	}
	resp, err := client.Do(req)

	if err != nil {
		g.log.Info("error making request")
		return models.GarantexRates{}, fmt.Errorf("error making request: %w", err)
	}
	var rates models.GarantexRates
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		g.log.Info("error decode response")
		return models.GarantexRates{}, fmt.Errorf("json decode error: %w", err)
	}
	if rates.Timestamp == 0 {
		g.log.Info("no rates found")
		return models.GarantexRates{}, errors.New("no rates found")
	}
	return rates, nil
}
