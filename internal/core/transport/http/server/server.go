package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Sergey-tech9087/golang-url-shortener/internal/core/logger"
	core_http_middleware "github.com/Sergey-tech9087/golang-url-shortener/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRoutes(routes ...*APIVersionRouter) {
	for _, router := range routes {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("Start HTTP server", zap.String("Addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("Listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("Shutdown HTTP server ...")

		shutDownContext, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutDownTimeout,
		)

		defer cancel()

		if err := server.Shutdown(shutDownContext); err != nil {
			_ = server.Close()

			return fmt.Errorf("Shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stoped")
	}

	return nil
}
