package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zhayt/transaction-service/config"
	"github.com/zhayt/transaction-service/transport/http/handler"
	"net"
	"time"
)

const _defaultShutdownTimeout = 3 * time.Second

type Server struct {
	handler         *handler.Handler
	App             *echo.Echo
	cfg             *config.Config
	Notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(handler *handler.Handler, cfg *config.Config) *Server {

	srv := &Server{
		handler:         handler,
		cfg:             cfg,
		shutdownTimeout: _defaultShutdownTimeout,
		Notify:          make(chan error, 1),
	}

	return srv
}

func (s *Server) StartHTTPServer() {
	s.App = s.buildingEngine()
	s.SetUpRoute()

	go func() {
		s.Notify <- s.App.Start(net.JoinHostPort("", s.cfg.AppPort))
		close(s.Notify)
	}()
}

func (s *Server) GracefullyShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.App.Shutdown(ctx)
}

func (s *Server) buildingEngine() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	return e
}
