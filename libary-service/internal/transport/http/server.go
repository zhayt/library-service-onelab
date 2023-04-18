package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/transport/http/handler"
	middleware2 "github.com/zhayt/user-storage-service/internal/transport/http/middleware"
	"time"
)

const (
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	App             *echo.Echo
	cfg             *config.Config
	handler         *handler.Handler
	mid             *middleware2.JWTAuth
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(cfg *config.Config, handler *handler.Handler, mid *middleware2.JWTAuth) *Server {
	srv := &Server{
		cfg:             cfg,
		handler:         handler,
		shutdownTimeout: _defaultShutdownTimeout,
		notify:          make(chan error, 1),
		mid:             mid,
	}

	return srv
}

func (s *Server) BuildingEngine() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	return e
}

func (s *Server) Start() {
	s.App = s.BuildingEngine()
	s.SetUpRoute()
	go func() {
		s.notify <- s.App.Start(fmt.Sprintf(":%s", s.cfg.AppPort))
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.App.Shutdown(ctx)
}
