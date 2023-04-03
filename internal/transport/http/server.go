package http

import (
	"context"
	"github.com/zhayt/user-storage-service/config"
	"net"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 10 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	Server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func NewServer(cfg *config.Config, router http.Handler) *Server {
	httpSrv := &http.Server{
		Addr: net.JoinHostPort("", cfg.Port),

		Handler: router,

		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	srv := &Server{
		Server:          httpSrv,
		shutdownTimeout: _defaultShutdownTimeout,
		notify:          make(chan error, 1),
	}

	return srv
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.Server.Shutdown(ctx)
}
