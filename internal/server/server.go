package server

import (
	"context"
	"fmt"
	"net/http"

	"indexer/internal/config"
	"indexer/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	srv    *http.Server
	cfg    config.HTTPConfig
}

func NewServer(cfg config.HTTPConfig) *Server {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	return &Server{
		router: router,
		cfg:    cfg,
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	middleware.Log.Infof("Starting HTTP server on %s", addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	middleware.Log.Info("Shutting down HTTP server...")
	return s.srv.Shutdown(ctx)
}
