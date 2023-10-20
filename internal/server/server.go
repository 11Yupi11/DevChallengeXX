package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"dev-challenge/db"
	"dev-challenge/internal/config"
	"dev-challenge/internal/handlers"
	"dev-challenge/internal/services"
)

type Server struct {
	conn    *sql.DB
	cfg     *config.Config
	log     logrus.FieldLogger
	storage db.Storage
}

func NewServer(conn *sql.DB, cfg *config.Config) (*Server, error) {
	s := &Server{
		conn:    conn,
		cfg:     cfg,
		storage: db.NewStorage(conn),
		log:     logrus.New(),
	}

	return s, nil
}

func (s *Server) Run() {
	router := chi.NewRouter()

	s.setupHandlers(router)

	srv := http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		s.log.Infof("start listening at %d", s.cfg.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// wait for terminate
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	<-terminate

	s.log.Info("server is shutting down...")

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		s.log.Errorf("could not gracefully shutdown the server: %v\n", err)
		logrus.Exit(1)
	}

	s.log.Info("server stopped")
	logrus.Exit(0)
}

func (s *Server) setupHandlers(router chi.Router) {
	router.Route("/api/v1", func(r chi.Router) {
		handler := handlers.ExcelLikeHandler{
			ELS: services.NewExcelLikeService(s.storage),
			Log: s.log,
		}
		handler.RegisterRoutes(r)
	})
}
