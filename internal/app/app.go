package app

import (
	"time"

	grpcapp "github.com/Sayfargo/url-shortener/internal/app/grpc"
	core_config "github.com/Sayfargo/url-shortener/internal/core/config"
	core_db "github.com/Sayfargo/url-shortener/internal/core/db/postgres"
	core_slogger "github.com/Sayfargo/url-shortener/internal/core/slogger"
	repository_urlshortener_postgres "github.com/Sayfargo/url-shortener/internal/repository/urlshortener/postgres"
	repository_urlshortener_redis "github.com/Sayfargo/url-shortener/internal/repository/urlshortener/redis"
	service_urlshortener "github.com/Sayfargo/url-shortener/internal/service/urlshortener"
	"github.com/redis/go-redis/v9"
)

type App struct {
	GRPCServer *grpcapp.App
}

func MustNew(cfg *core_config.Config) *App {

	log := core_slogger.MustNew(cfg.Logger.Level, cfg.Logger.Dir)

	pool, err := core_db.InitWithRerty(cfg, log, 5, time.Second*2)
	if err != nil {
		panic(err.Error())
	}
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	postgresRepo := repository_urlshortener_postgres.NewPostgresRepository(pool)
	redisRepo := repository_urlshortener_redis.NewRedisRepository(client)

	service := service_urlshortener.NewUrlShortenerService(postgresRepo, redisRepo, log)

	grcpApp := grpcapp.NewApp(log, service, cfg.GRPC.Port)

	return &App{
		GRPCServer: grcpApp,
	}
}
