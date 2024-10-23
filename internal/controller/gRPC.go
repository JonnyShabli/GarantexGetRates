package controller

import (
	"context"
	"fmt"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
	"github.com/JonnyShabli/GarantexGetRates/internal/repository"
	"go.opentelemetry.io/otel/trace"

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
	tracer  trace.Tracer
	pb.GgrServer
}

func NewGRPCObj(log *zap.Logger, repo repository.GgrRepoInterface, tracer trace.Tracer) *GRPCObj {
	return &GRPCObj{
		Service: service.NewGgrService(log),
		log:     log,
		Repo:    repo,
		tracer:  tracer,
	}
}

func (g *GRPCObj) GarantexGetRates(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	_, span := g.tracer.Start(ctx, "GetRates")
	res, err := g.Service.GetRates(ctx, req.GetPair())
	span.End()
	if err != nil {
		return nil, fmt.Errorf("get rates: %w", err)
	}

	dto := models.RatesToDB{
		Timestamp: res.Timestamp,
		Ask:       res.Asks[0],
		Bid:       res.Bids[0],
	}

	_, span = g.tracer.Start(ctx, "InsertRates")
	err = g.Repo.InsertRates(ctx, dto)
	span.End()
	if err != nil {
		g.log.Error("insert rates", zap.Error(err))
		return nil, fmt.Errorf("insert rates: \"%w\"", err)
	}

	return &pb.Response{Msg: fmt.Sprintf("Succesfuly safe to DB rates for %v", req.GetPair())}, nil
}
