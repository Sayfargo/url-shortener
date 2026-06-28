package service_urlshortener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	core_slogger "github.com/Sayfargo/url-shortener/internal/core/slogger"
	repository_urlshortener_postgres "github.com/Sayfargo/url-shortener/internal/repository/urlshortener/postgres"
	repository_urlshortener_redis "github.com/Sayfargo/url-shortener/internal/repository/urlshortener/redis"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type UrlPostgresRepository interface {
	Create(ctx context.Context, originalURL, shortCode string) error
	Get(ctx context.Context, shortCode string) (string, error)
}

type UrlRedisRepository interface {
	Set(ctx context.Context, originalURL, shortCode string) error
	Get(ctx context.Context, shortCode string) (string, error)
}

type UrlShortenerService struct {
	postgresRepo UrlPostgresRepository
	redisRepo    UrlRedisRepository
	log          *core_slogger.Slogger
}

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	ErrUrlDoesNotExists = errors.New("url doesn't exists")
	ErrReponseTimeout   = errors.New("response timeout")
	ErrCodeAlredyExists = errors.New("short code already exists")
)

func NewUrlShortenerService(
	postgresRepo UrlPostgresRepository,
	redisRepo UrlRedisRepository,
	log *core_slogger.Slogger,
) *UrlShortenerService {
	return &UrlShortenerService{
		postgresRepo: postgresRepo,
		redisRepo:    redisRepo,
		log:          log,
	}
}

func (u *UrlShortenerService) Get(ctx context.Context, shortCode string) (string, error) {

	const op = "service_urlshortener.Get"

	l := u.log.With(
		slog.String("op", op),
	)

	url, err := u.redisRepo.Get(ctx, shortCode)

	if err == nil {
		l.Debug(
			"URL retrieved via Redis",
			slog.String("url", url),
		)
		return url, nil
	}

	if !errors.Is(err, repository_urlshortener_redis.ErrNotExists) {
		l.Debug(
			"Redis error. Falling back to Postgres",
			slog.String("err", err.Error()),
		)
	}

	url, err = u.postgresRepo.Get(ctx, shortCode)
	if err != nil {
		if errors.Is(err, repository_urlshortener_redis.ErrNotExists) {
			l.Debug(
				"failed to find URL",
				slog.String("err", err.Error()),
			)
			return "", ErrUrlDoesNotExists
		}

		l.Error(
			"failed to get URL via Postgres",
			slog.String("err", err.Error()),
		)

		return "", fmt.Errorf("postgres error : %w", err)
	}

	_ = u.redisRepo.Set(ctx, url, shortCode)

	l.Debug(
		"URL retrieved via Postgres",
		slog.String("url", url),
	)

	return url, nil

}

func (u *UrlShortenerService) CreateShortCode(ctx context.Context, url string) (string, error) {

	const op = "service_urlshortener.CreateShorted"

	l := u.log.With(
		slog.String("op", op),
		slog.String("original_url", url),
	)

	var (
		shortCode string
		err       error
	)

	// if conflict do retry
	for attempt := 0; attempt < 3; attempt++ {

		shortCode, err = gonanoid.Generate(alphabet, 8)
		if err != nil {
			l.Error(
				"failed to generate short code",
				slog.String("err", err.Error()),
			)

			return "", fmt.Errorf("failed to generate short code : %w", err)
		}

		if err := u.postgresRepo.Create(ctx, url, shortCode); err != nil {

			if errors.Is(err, repository_urlshortener_postgres.ErrAlredyExists) {
				l.Warn(
					"short code conflict detected, retrying...",
					slog.String("short_code", shortCode),
					slog.Int("attempt", attempt+1),
				)

				if attempt == 2 {
					return "", repository_urlshortener_postgres.ErrAlredyExists
				}

				continue
			}

			l.Error(
				"failed to create short code via Postgres",
				slog.String("err", err.Error()),
			)

			return "", fmt.Errorf("postgres error : %w", err)
		}

		break

	}

	if err := u.redisRepo.Set(ctx, url, shortCode); err != nil {
		l.Error(
			"failed to set URL and short code via Redis",
			slog.String("short_code", shortCode),
			slog.String("err", err.Error()),
		)
	}

	l.Debug(
		"Short code succesfully created",
		slog.String("url", url),
	)

	return shortCode, nil

}
