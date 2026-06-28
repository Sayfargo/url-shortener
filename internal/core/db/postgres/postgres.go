package core_db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	core_config "github.com/Sayfargo/url-shortener/internal/core/config"
	core_slogger "github.com/Sayfargo/url-shortener/internal/core/slogger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitWithRerty(
	cfg *core_config.Config,
	log *core_slogger.Slogger,
	attempts int,
	delay time.Duration,
) (*pgxpool.Pool, error) {

	log.Info(
		"initialization database",
	)

	pgxCfg, err := pgxpool.ParseConfig(cfg.Postgres.URL)
	if err != nil {
		log.Error(
			"failed to parse database config",
			slog.String("err", err.Error()),
		)

		return nil, fmt.Errorf("parse config : %w", err)
	}

	pgxCfg.MaxConns = cfg.Postgres.MaxConns
	pgxCfg.MinConns = cfg.Postgres.MinConns
	pgxCfg.HealthCheckPeriod = cfg.Postgres.HealthCheckPeriod

	var (
		try            int = 0
		connectionErrs error
	)

	for attempt := 0; attempt < attempts; attempt++ {

		try = attempt + 1

		if attempt > 0 {
			time.Sleep(delay)
			delay *= 2
		}

		log.Info(
			"try to connect to database",
			slog.Int("attempt", try),
			slog.Int("attempts left", attempts-try),
		)

		pool, err := pgxpool.NewWithConfig(context.Background(), pgxCfg)
		if err != nil {
			log.Warn(
				"failed to create database with pool",
				slog.String("err", err.Error()),
			)
			connectionErrs = errors.Join(err)
			continue
		}

		if err := pool.Ping(context.Background()); err != nil {
			log.Warn(
				"failed to ping database",
				slog.String("err", err.Error()),
			)
			connectionErrs = errors.Join(err)
			pool.Close()
			continue
		}

		log.Info(
			"database initialized successfully",
		)

		return pool, nil

	}

	log.Error(
		"failed to connect to database after all attempts",
		slog.Int("max_conns", int(cfg.Postgres.MaxConns)),
		slog.Int("min_conns", int(cfg.Postgres.MinConns)),
		slog.Duration("health_check_period", cfg.Postgres.HealthCheckPeriod),
		slog.String("errors", connectionErrs.Error()),
	)

	return nil, fmt.Errorf("failed to connect db : %w", connectionErrs)

}
