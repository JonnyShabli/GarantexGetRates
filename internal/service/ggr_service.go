package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"go.uber.org/zap"
	"net/http"
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		g.log.Info("error forming request")
		return models.GarantexRates{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		g.log.Info("error making request")
		return models.GarantexRates{}, err
	}
	var rates models.GarantexRates
	err = json.NewDecoder(resp.Body).Decode(&rates)
	if err != nil {
		g.log.Info("error decode response")
		return models.GarantexRates{}, err
	}

	return rates, nil
}
