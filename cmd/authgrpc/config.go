package main

import (
	"errors"
	"github.com/vaberof/auth-grpc/internal/domain/auth"
	"github.com/vaberof/auth-grpc/pkg/config"
	"github.com/vaberof/auth-grpc/pkg/database/postgres"
	"github.com/vaberof/auth-grpc/pkg/database/redis"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcclient"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcserver"
	"os"
)

type AppConfig struct {
	Server      grpcserver.ServerConfig
	AuthService auth.Config
	Postgres    postgres.Config
	Redis       redis.Config

	NotificationService grpcclient.NotificationServiceClientConfig
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig grpcserver.ServerConfig
	err := config.ParseConfig(provider, "app.grpc.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var authConfig auth.Config
	err = config.ParseConfig(provider, "app.auth-service", &authConfig)
	if err != nil {
		return nil, err
	}

	var postgresConfig postgres.Config
	err = config.ParseConfig(provider, "app.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}
	postgresConfig.User = os.Getenv("POSTGRES_USER")
	postgresConfig.Password = os.Getenv("POSTGRES_PASSWORD")

	var redisConfig redis.Config
	err = config.ParseConfig(provider, "app.redis", &redisConfig)
	if err != nil {
		return nil, err
	}

	var notificationServiceConfig grpcclient.NotificationServiceClientConfig
	err = config.ParseConfig(provider, "app.grpc.client.notification-service", &notificationServiceConfig)
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		Server:              serverConfig,
		AuthService:         authConfig,
		Postgres:            postgresConfig,
		Redis:               redisConfig,
		NotificationService: notificationServiceConfig,
	}

	return &appConfig, nil
}
