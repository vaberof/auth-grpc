package grpcserver

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type AppServer struct {
	Server  *grpc.Server
	config  *ServerConfig
	logger  *slog.Logger
	address string
}

func New(config *ServerConfig, logs *logs.Logs) *AppServer {
	logger := logs.WithName("gRPC-server")

	loggingOpts := getLoggingOpts()
	recoveryOpts := getRecoveryOpts(logger)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logs.GetLogger()), loggingOpts...),
	))

	appServer := &AppServer{
		Server:  grpcServer,
		config:  config,
		logger:  logger,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}

	return appServer
}

func (server *AppServer) StartAsync() <-chan error {
	server.logger.Info("Starting gRPC server")

	exitChannel := make(chan error, 1)

	listener, err := net.Listen("tcp", server.address)
	if err != nil {
		exitChannel <- err
		return exitChannel
	}

	go func() {
		err = server.Server.Serve(listener)
		if err != nil {
			server.logger.Error("Failed to start gRPC server", slog.Group("error", err))

			exitChannel <- err
		} else {
			exitChannel <- nil
		}
	}()

	server.logger.Info("Started gRPC server", slog.Group("gRPC-server", "address", server.address))

	return exitChannel
}

func (server *AppServer) GracefulShutdown() {
	server.logger.Info("Stopping gRPC server")

	server.Server.GracefulStop()

	server.logger.Info("gRPC server is stopped")
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func getLoggingOpts() []logging.Option {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	return loggingOpts
}

func getRecoveryOpts(logger *slog.Logger) []recovery.Option {
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	return recoveryOpts
}
