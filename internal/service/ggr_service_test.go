package service

import (
	"context"
	"testing"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetRates(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	defer func() { _ = logger.Sync() }()

	type args struct {
		ctx  context.Context
		pair string
	}
	tests := []struct {
		name   string
		args   args
		expect models.GarantexRates
	}{
		{
			name: "GetRates_NoError",
			args: args{
				ctx:  context.Background(),
				pair: "btcrub",
			},
		},
		{
			name: "GetRates_BadPair",
			args: args{
				ctx:  context.Background(),
				pair: "BADPAIR",
			},
			expect: models.GarantexRates{
				Timestamp: 0,
				Asks:      nil,
				Bids:      nil,
			},
		},
	}

	for _, tt := range tests {
		servObj := NewGgrService(logger)
		t.Run(tt.name, func(t *testing.T) {
			_, err := servObj.GetRates(tt.args.ctx, tt.args.pair)
			require.NoError(t, err)
		})

	}
}

func TestGetRatesError(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	defer func() { _ = logger.Sync() }()

	type args struct {
		ctx  context.Context
		pair string
	}
	tests := []struct {
		name   string
		args   args
		expect models.GarantexRates
	}{
		{
			name: "GetRates_badPair",
			args: args{
				ctx:  context.Background(),
				pair: "BADPAIR",
			},
			expect: models.GarantexRates{
				Timestamp: 0,
				Asks:      nil,
				Bids:      nil,
			},
		},
	}

	for _, tt := range tests {
		servObj := NewGgrService(logger)
		t.Run(tt.name, func(t *testing.T) {
			got, err := servObj.GetRates(tt.args.ctx, tt.args.pair)
			if assert.NoError(t, err, "expecting no errors") {
				require.Equal(t, tt.expect, got)
			}
		})

	}
}
