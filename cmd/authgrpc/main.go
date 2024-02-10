package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/auth-grpc/internal/app/entrypoint/grpc/auth"
	authservice "github.com/vaberof/auth-grpc/internal/domain/auth"
	userservice "github.com/vaberof/auth-grpc/internal/domain/user"
	"github.com/vaberof/auth-grpc/internal/infra/integration/grpc/notificationservice"
	"github.com/vaberof/auth-grpc/internal/infra/storage/postgres/pguser"
	redisstorage "github.com/vaberof/auth-grpc/internal/infra/storage/redis"
	"github.com/vaberof/auth-grpc/pkg/database/postgres"
	"github.com/vaberof/auth-grpc/pkg/database/redis"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcclient"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcserver"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()

	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	logger := logs.New(os.Stdout, nil)

	appConfig := mustGetAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	postgresManagedDb, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	redisManagedDb, err := redis.New(&appConfig.Redis)
	if err != nil {
		panic(err)
	}

	notificationServiceGrpcClient, err := grpcclient.New(&appConfig.NotificationService)
	if err != nil {
		panic(err)
	}

	redisStorage := redisstorage.NewRedisStorage(redisManagedDb.RedisDb)
	pgUserStorage := pguser.NewPgUserStorage(postgresManagedDb.PostgresDb)

	notificationService := notificationservice.New(notificationServiceGrpcClient, logger)
	userService := userservice.NewUserService(pgUserStorage, logger)

	authService := authservice.NewAuthService(&appConfig.AuthService, userService, notificationService, redisStorage, logger)

	// TODO: implement general graceful shutdown for databases and server

	grpcServer := grpcserver.New(&appConfig.Server, logger)

	auth.Register(grpcServer.Server, authService)

	grpcServerErrorCh := grpcServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signalValue := <-quitCh:
		logger.GetLogger().Info("stopping application", slog.String("signal", signalValue.String()))

		grpcServer.Shutdown()
	case err = <-grpcServerErrorCh:
		logger.GetLogger().Info("stopping application", slog.String("gRPC server error", err.Error()))

		grpcServer.Shutdown()
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
