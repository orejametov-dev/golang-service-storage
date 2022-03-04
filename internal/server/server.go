package server

import (
	"context"
	"fmt"
	"net/http"
	"orejametov/service-storage/internal/config"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			ReadTimeout:    time.Second * cfg.HTTP.ReadTimeout,
			WriteTimeout:   time.Second * cfg.HTTP.ReadTimeout,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run() error {
	fmt.Println("Service is runned successfully!!!")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
