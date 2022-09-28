package server

import (
	"database/sql"
	"net/http"
	"time"
	
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer http.Server
	config *Config
	logger *logrus.Logger
}

func NewServer(r http.Handler, config *Config) *Server {
	logrus.New()
	return &Server{
		httpServer : http.Server{
			Addr:              config.ServerAddr,
			Handler:           r,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
		},
		config: config,
		logger: logrus.New(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres password=postgres")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
