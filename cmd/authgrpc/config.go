package main

import (
	"errors"
	"github.com/vaberof/auth-grpc/pkg/config"
	"github.com/vaberof/auth-grpc/pkg/database/postgres"
	"github.com/vaberof/auth-grpc/pkg/database/redis"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcserver"
)

type AppConfig struct {
	Server   grpcserver.ServerConfig
	Postgres postgres.Config
	Redis    redis.Config
}

func getAppConfig(sources ...string) AppConfig {
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

	var postgresConfig postgres.Config
	err = config.ParseConfig(provider, "app.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}

	var redisConfig redis.Config
	err = config.ParseConfig(provider, "app.redis", &redisConfig)
	if err != nil {
		return nil, err
	}

	config := AppConfig{
		Server:   serverConfig,
		Postgres: postgresConfig,
		Redis:    redisConfig,
	}

	return &config, nil
}
