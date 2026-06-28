package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	core_slogger "github.com/Sayfargo/url-shortener/internal/core/slogger"
	service_urlshortener "github.com/Sayfargo/url-shortener/internal/service/urlshortener"
	transport_grpc_urlshortener "github.com/Sayfargo/url-shortener/internal/transport/grpc/urlshortener"
	"google.golang.org/grpc"
)

type App struct {
	log        *core_slogger.Slogger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(
	log *core_slogger.Slogger,
	svc *service_urlshortener.UrlShortenerService,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	transport_grpc_urlshortener.RegisterServer(gRPCServer, svc, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		l.Error(
			"failed to listen",
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("failed to listen : %w", err)
	}

	l.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		l.Error(
			"failed to server",
			slog.String("err", err.Error()),
		)
		return fmt.Errorf("failed to server : %w", err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

	a.log.Close()

}
