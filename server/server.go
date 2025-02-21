package server

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"path"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	server http.Server
	logger zerolog.Logger
}

var (
	//go:embed static/*
	staticFS embed.FS
)

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix, filepath string) bool {
	f, err := e.Open(path.Join(prefix, filepath))
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func New(cfg *Config) (*Server, error) {

	// Throw an error if no username / password was provided
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("username and password must be supplied")
	}

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

	// Static files
	r.Use(
		static.Serve("/", embedFileSystem{http.FS(staticFS)}),
	)

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
