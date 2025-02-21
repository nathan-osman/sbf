package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server http.Server
	logger zerolog.Logger
}

func New(cfg *Config) (*Server, error) {

	// Switch Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	var (
		r = gin.New()
		s = &Server{
			server: http.Server{
				Addr:    cfg.Addr,
				Handler: r,
			},
			logger: log.With().Str("package", "server").Logger(),
		}
	)

	// Admin pages
	admin := r.Group("/admin")
	admin.Use(
		gin.BasicAuth(
			gin.Accounts{cfg.Username: cfg.Password},
		),
	)

	admin.GET("/", s.adminIndex)

	// The home page
	r.GET("/", s.index)

	// Start the goroutine that listens for incoming connections
	go func() {
		defer s.logger.Info().Msg("server stopped")
		s.logger.Info().Msg("server started")
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Msg(err.Error())
		}
	}()

	return s, nil
}

// Close shuts down the server.
func (s *Server) Close() {
	s.server.Shutdown(context.Background())
}
