package controller

import (
	"context"
	"testing"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	pb "github.com/JonnyShabli/GarantexGetRates/internal/proto/ggr"
	"github.com/JonnyShabli/GarantexGetRates/internal/repository/mock_repo"
	"github.com/JonnyShabli/GarantexGetRates/pkg/tracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetRates(t *testing.T) {
	type args struct {
		ctx context.Context
		pb.Request
	}

	tests := []struct {
		name string
		args args
		want models.RatesToDB
	}{
		{
			name: "Get and Save no error",
			args: args{
				ctx: context.Background(),
				Request: pb.Request{
					Pair: "btcrub",
				},
			},
		},
	}

	for _, tt := range tests {
		logger, _ := zap.NewDevelopment()
		repo := mock_repo.NewMockRepo()
		traceMgr, _ := tracer.InitTracer("", "")
		grpcObj := NewGRPCObj(logger, repo, traceMgr)

		t.Run(tt.name, func(t *testing.T) {
			_, err := grpcObj.GetRates(tt.args.ctx, &tt.args.Request)
			if assert.NoError(t, err) {
				//assert.Equal(t, tt.want, mock_repo.GetRates(repo, 1))
				require.NotEmpty(t, mock_repo.GetRates(repo, repo.Idx))
			}
		})
	}
}
