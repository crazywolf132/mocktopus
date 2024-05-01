package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	server     *http.Server
	mu         sync.Mutex
	Router     *http.ServeMux
	logger     *zap.Logger
	middleware []Middleware
}

type Handler interface {
	Handle(http.ResponseWriter, *http.Request) bool
}

type Middleware func(http.Handler) http.Handler

func CreateServer(config *viper.Viper) (*Server, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	router := http.NewServeMux()

	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", config.GetInt("server.port")),
		},
		Router: router,
		logger: logger,
	}, nil
}

func (s *Server) Start() error {
	s.server.Handler = s.Router

	s.logger.Info("Starting mock server", zap.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start mock server: %v", err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("Stopping mock server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop mock server: %v", err)
	}
	return nil
}
