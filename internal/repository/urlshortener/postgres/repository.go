package repository_urlshortener_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

var (
	ErrNotExists    = errors.New("record doesn't exists")
	ErrAlredyExists = errors.New("already exists")
)

func (p *PostgresRepository) Create(ctx context.Context, originalURL, shortCode string) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}

	query := `INSERT INTO shorted_urls (original_url, short_code) VALUES ($1, $2)`

	_, err := p.pool.Exec(ctx, query, originalURL, shortCode)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return ErrAlredyExists
			}
		}

		return fmt.Errorf("failed to exec : %w", err)
	}

	return nil

}

func (p *PostgresRepository) Get(ctx context.Context, shortCode string) (string, error) {

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	query := `SELECT original_url FROM shorted_urls WHERE short_code=$1`

	var url string

	err := p.pool.QueryRow(ctx, query, shortCode).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotExists
		}
		return "", fmt.Errorf("failed to query row : %w", err)
	}

	return url, nil
}
