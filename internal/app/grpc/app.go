package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/ArtShib/urlshortener/internal/grpc/interceptors"
	"github.com/ArtShib/urlshortener/internal/grpc/server"
	"github.com/ArtShib/urlshortener/internal/lib/auth"
	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/service"
	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int, authSvc *auth.Service, urlSvc *service.URLService) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthInterceptor(authSvc),
		),
	)
	server.Register(gRPCServer, urlSvc)
	return &App{
		logger:     logger,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (a *App) Start(ctx context.Context) error {
	logHelper := loghelper.New(a.logger, "ServerGRPC.Start")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return logHelper.LogAndReturnError(ctx, "Error listen port", err)
	}
	logHelper.LogDebug(ctx, "Server GRPC start")
	if err := a.gRPCServer.Serve(l); err != nil {
		return logHelper.LogAndReturnError(ctx, "Error starting server", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	logHelper := loghelper.New(a.logger, "ServerGRPC.Stop")
	a.gRPCServer.GracefulStop()
	logHelper.LogDebug(ctx, "Server GRPC stop")
}
