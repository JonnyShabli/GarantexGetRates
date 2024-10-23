package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/JonnyShabli/GarantexGetRates/internal/controller"
	"github.com/JonnyShabli/GarantexGetRates/internal/db"
	"github.com/JonnyShabli/GarantexGetRates/internal/pkg/health"
	pkghttp "github.com/JonnyShabli/GarantexGetRates/internal/pkg/http"
	"github.com/JonnyShabli/GarantexGetRates/internal/pkg/sig"
	"github.com/JonnyShabli/GarantexGetRates/internal/pkg/tracer"
	pb "github.com/JonnyShabli/GarantexGetRates/internal/proto/ggr"
	"github.com/JonnyShabli/GarantexGetRates/internal/repository"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// читаем настройки из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// читаем флаги запуска
	var dbConnString string
	flag.StringVar(&dbConnString, "db", "", "Database connection string")
	flag.Parse()

	// Создаем errgroup и контекст
	g, ctx := errgroup.WithContext(context.Background())

	// Создаем логгер
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("fail to make logger", err)
	}
	defer func() { _ = logger.Sync() }()

	// Создаем OpenTelemetry tracer
	traceMgr, err := tracer.InitTracer(os.Getenv("TRACE_ADDR"), "GGR Service")
	if err != nil {
		logger.Fatal("fail to init tracer", zap.Error(err))
	}

	// создаем объект слоя репозиторий SQL DB
	if dbConnString == "" {
		dbConnString = fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"))
	}
	dbConn := db.NewConn(dbConnString)
	storage := repository.NewGgrRepo(logger, dbConn)

	// запускаем миграции БД
	err = db.InitDB(dbConnString)
	if err != nil {
		logger.Fatal("fail to init db migrations", zap.Error(err))
	}

	// создаем объект слоя controller - gRPC хэндлер
	logger.Info("creating grpcHandler object")
	grpcHandler := controller.NewGRPCObj(logger, storage, traceMgr)

	// создаем и настраиваем gRPC сервер
	logger.Info("creating grpcSever object")
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
		),
	)

	// регистрируем методы в gRPC сервере
	pb.RegisterGgrServer(grpcServer, grpcHandler)

	// для тестирования через Postman
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	g.Go(func() error {
		connString := fmt.Sprintf("%v:%v", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))
		logger.Info(fmt.Sprintf("Listen on: %v", connString))

		listen, err := net.Listen("tcp", connString)
		if err != nil {
			logger.Fatal("error listen :"+connString, zap.String("error", err.Error()))
		}

		errListen := make(chan error, 1)
		go func() {
			errListen <- grpcServer.Serve(listen)
		}()

		select {
		case <-ctx.Done():
			grpcServer.GracefulStop()
			return nil
		case err = <-errListen:
			return fmt.Errorf("can't run grpc server: %w", err)
		}
	})

	// Ждем сигналы ОС для завершения работы
	g.Go(func() error { return sig.ListenSignal(ctx, logger) })

	// metrics.
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"go_project": "ggr"}, registry)
	registerer.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	techHandler := pkghttp.NewHandler("/", pkghttp.DefaultTechOptions(registry))

	// Запускаем http сервер
	g.Go(func() error {
		return pkghttp.RunServer(ctx, os.Getenv("HTTP_PRIVATE_ADDR"), logger, techHandler)
	})

	health.SetStatus(http.StatusOK)
	logger.Info("Waiting for signal")
	err = g.Wait()
	if err != nil && !errors.Is(err, sig.ErrSignalReceived) {
		logger.With(zap.String("Exit reason", err.Error()))
	}
}
