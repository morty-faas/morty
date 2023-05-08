package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/morty-faas/morty/registry/builder"
	"github.com/morty-faas/morty/registry/config"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morty-faas/morty/registry/storage"
	"github.com/morty-faas/morty/registry/storage/s3"
)

const (
	healthEndpoint = "/healthz"
)

type (
	Server struct {
		cfg     *config.Config
		storage storage.Storage
		builder *builder.Builder
	}

	APIError struct {
		StatusCode int    `json:"status"`
		Message    string `json:"message"`
	}
)

func NewServer() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	s3, err := s3.NewStorage(cfg.Storage.S3)
	if err != nil {
		return nil, err
	}

	bld, err := builder.NewBuilder()
	if err != nil {
		return nil, err
	}

	server := &Server{
		cfg:     cfg,
		storage: s3,
		builder: bld,
	}

	return server, nil
}

func (s *Server) Serve() {
	p := s.cfg.Port

	ctx, stop := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", p),
		Handler: s.router(),
	}

	// Listen for syscall signals for process to interrupt/quit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				cancel()
				log.Fatal("graceful shutdown timed out... forcing exit")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		stop()
	}()

	log.Printf("function registry is listening on 0.0.0.0:%d\n", p)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-ctx.Done()
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)

	r.Post("/v1/functions/build", s.BuildHandler)
	r.Get("/v1/functions/{id}/{version}", s.GetImageHandler)
	r.Get(healthEndpoint, s.HealthcheckHandler)

	return r
}

func (s *Server) JSONResponse(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) APIErrorResponse(w http.ResponseWriter, err *APIError) {
	log.Error(err.Message)
	s.JSONResponse(w, err.StatusCode, err)
}
