package transport_grpc_urlshortener

import (
	"context"
	"errors"
	"log/slog"

	urlv1 "github.com/Sayfargo/protos/gen/go/url"
	core_slogger "github.com/Sayfargo/url-shortener/internal/core/slogger"
	repository_urlshortener_postgres "github.com/Sayfargo/url-shortener/internal/repository/urlshortener/postgres"
	service_urlshortener "github.com/Sayfargo/url-shortener/internal/service/urlshortener"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UrlShortenerService interface {
	CreateShortCode(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, shortCode string) (string, error)
}

type serverAPI struct {
	urlv1.UnimplementedUrlShortenerServer

	svc UrlShortenerService
	log *core_slogger.Slogger
	val *validator.Validate
}

func RegisterServer(
	gRPC *grpc.Server,
	svc UrlShortenerService,
	log *core_slogger.Slogger,
) {
	urlv1.RegisterUrlShortenerServer(gRPC, &serverAPI{svc: svc, log: log, val: validator.New()})

}

func (s *serverAPI) Get(ctx context.Context, request *urlv1.GetUrlRequest) (*urlv1.GetUrlResponse, error) {
	const pen = "Get"

	l := s.log.With(
		slog.String("pen", pen),
	)

	if request.ShortCode == "" {
		l.Debug(
			"short code is empty",
		)

		return nil, status.Error(
			codes.InvalidArgument,
			"invalid short code",
		)

	}

	url, err := s.svc.Get(ctx, request.ShortCode)
	if err != nil {
		if errors.Is(err, service_urlshortener.ErrUrlDoesNotExists) {

			l.Debug(
				"url not found",
				slog.String("err", err.Error()),
			)

			return nil, status.Error(
				codes.NotFound,
				"url not found",
			)
		}

		l.Error(
			"failed to get url",
			slog.String("err", err.Error()),
		)

		return nil, status.Error(
			codes.Internal,
			"internal error",
		)

	}

	return &urlv1.GetUrlResponse{
		Url: url,
	}, nil

}

func (s *serverAPI) Create(ctx context.Context, request *urlv1.CreateShortCodeRequest) (*urlv1.CreateShortCodeResponse, error) {
	const pen = "Short"

	l := s.log.With(
		slog.String("pen", pen),
	)

	if err := s.val.Var(request.Url, "url"); err != nil {
		l.Debug(
			"failed to validate url",
			slog.String("err", err.Error()),
		)

		return nil, status.Error(
			codes.InvalidArgument,
			"invalid url",
		)
	}

	shortCode, err := s.svc.CreateShortCode(ctx, request.Url)
	if err != nil {
		if errors.Is(err, repository_urlshortener_postgres.ErrAlredyExists) {
			return nil, status.Error(codes.AlreadyExists, "try again O_o")
		}
		return nil, status.Error(
			codes.Internal,
			"internal error",
		)
	}

	return &urlv1.CreateShortCodeResponse{
		ShortCode: shortCode,
	}, nil

}
