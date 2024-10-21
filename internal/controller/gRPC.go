package controller

import (
	"context"
	"fmt"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"github.com/JonnyShabli/GarantexGetRates/internal/repository"

	pb "github.com/JonnyShabli/GarantexGetRates/internal/proto/ggr"
	"github.com/JonnyShabli/GarantexGetRates/internal/service"
	"go.uber.org/zap"
)

type GRPCInterface interface {
	GarantexGetRates(ctx context.Context, req *pb.Request) (*pb.Response, error)
}

type GRPCObj struct {
	Service service.GgrServiceInterface
	log     *zap.Logger
	Repo    repository.GgrRepoInterface
	pb.GgrServer
}

func NewGRPCObj(log *zap.Logger, repo repository.GgrRepoInterface) *GRPCObj {
	return &GRPCObj{
		Service: service.NewGgrService(log),
		log:     log,
		Repo:    repo,
	}
}

func (g *GRPCObj) GarantexGetRates(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	res, err := g.Service.GetRates(ctx, req.GetPair())
	if err != nil {
		return nil, fmt.Errorf("get rates: %w", err)
	}

	dto := models.RatesToDB{
		Timestamp: res.Timestamp,
		Ask:       res.Asks[0],
		Bid:       res.Bids[0],
	}

	err = g.Repo.InsertRates(ctx, dto)
	if err != nil {
		g.log.Error("insert rates", zap.Error(err))
		return nil, fmt.Errorf("insert rates: \"%w\"", err)
	}

	return &pb.Response{Msg: fmt.Sprintf("Succesfuly safe to DB rates for %v", req.GetPair())}, nil
}
