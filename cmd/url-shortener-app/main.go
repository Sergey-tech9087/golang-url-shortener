package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Sergey-tech9087/golang-url-shortener/internal/core/logger"
	core_pgx_pool "github.com/Sergey-tech9087/golang-url-shortener/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/middleware"
	core_http_server "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/server"
	urls_inmemory_repository "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/repository/inmemory"
	urls_postgres_repository "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/repository/postgres"
	urls_service "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/service"
	url_transport_http "github.com/Sergey-tech9087/golang-url-shortener/internal/features/shortener/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "postgres"
	}

	var urlsRepository urls_service.UrlsRepository

	switch storageType {
	case "postgres":
		logger.Debug("Initializing postgres connection pool")
		pool, err := core_pgx_pool.NewPool(
			ctx,
			core_pgx_pool.NewConfigMust(),
		)
		if err != nil {
			logger.Fatal("failed to init postgres connection pool", zap.Error(err))
		}
		defer pool.Close()
		urlsRepository = urls_postgres_repository.NewUrlsRepository(pool)

	case "inmemory":
		logger.Debug("Initializing in-memory storage")
		urlsRepository = urls_inmemory_repository.NewUrlsRepository()

	default:
		logger.Fatal("unknown storage type", zap.String("storage_type", storageType))
	}

	logger.Debug("Initializing feature", zap.String("feature", "shortener_url"))
	urlService := urls_service.NewUrlsService(urlsRepository)
	urlsTransportHTTP := url_transport_http.NewUrlsHTTPHandler(urlService)

	logger.Debug("Initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(urlsTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
