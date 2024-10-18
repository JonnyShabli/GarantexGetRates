package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/JonnyShabli/GarantexGetRates/internal/controller"
	"github.com/JonnyShabli/GarantexGetRates/internal/db"
	pb "github.com/JonnyShabli/GarantexGetRates/internal/proto/ggr"
	"github.com/JonnyShabli/GarantexGetRates/internal/repository"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
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
	g, _ := errgroup.WithContext(context.Background())

	// Создаем логгер
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("fail to make logger", err)
	}
	defer func() { _ = logger.Sync() }()

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

	// создаем объект слоя controller - gRPC хэндлер
	logger.Info("creating grpcHandler object")
	grpcHandler := controller.NewGRPCObj(logger, storage)

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

		return grpcServer.Serve(listen)
	})

	g.Wait()
}
