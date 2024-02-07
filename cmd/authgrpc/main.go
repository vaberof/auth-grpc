package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/auth-grpc/internal/app/entrypoint/grpc/auth"
	"github.com/vaberof/auth-grpc/pkg/database/postgres"
	"github.com/vaberof/auth-grpc/pkg/database/redis"
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

	appConfig := getAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	_, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	_, err = redis.New(&appConfig.Redis)
	if err != nil {
		panic(err)
	}

	// TODO: init service

	grpcServer := grpcserver.New(&appConfig.Server, logger)

	auth.Register(grpcServer.Server)

	grpcServerErrorCh := grpcServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-quitCh:
		logger.GetLogger().Info("stopping application", slog.String("signal", sign.String()))

		grpcServer.GracefulShutdown()
	case err := <-grpcServerErrorCh:
		logger.GetLogger().Info("stopping application", slog.String("gRPC server error", err.Error()))

		grpcServer.GracefulShutdown()
	}

}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
