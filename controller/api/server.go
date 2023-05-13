package api

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/api/handlers"
	"github.com/morty-faas/morty/controller/config"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/types"
	"github.com/morty-faas/morty/pkg/telemetry"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type server struct {
	cfg   *config.Config
	state state.State
	orch  orchestration.Orchestrator
}

// New initializes a new API server.
// The method will return an error if the bootstrap process encounters one.
func New(cfg *config.Config) (*server, error) {
	orch, err := cfg.OrchestratorFactory()
	if err != nil {
		return nil, err
	}

	state, err := cfg.StateFactory(func(instance string) {
		if err := orch.DeleteFunctionInstance(context.Background(), &types.FnInstance{Id: instance}); err != nil {
			log.Errorf("Failed to automatically delete instance %s: %v", instance, err)
		}
	})
	if err != nil {
		return nil, err
	}

	srv := &server{cfg, state, orch}
	srv.getInitialState()
	return srv, nil
}

// ListenAndServe starts the server on the configured port and configure everything
// needed to support the graceful shutdown.
func (s *server) ListenAndServe() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.makeRouter(),
	}

	// Run the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logrus.Infof("API server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Could not start internal HTTP server: %v", err)
		}
	}()

	metricsAddr := fmt.Sprintf(":%d", s.cfg.MetricsPort)
	metricsSrv := telemetry.NewMetricsServer(metricsAddr)

	go func() {
		logrus.Infof("Metrics server listening on %s", metricsAddr)
		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Could not start Metrics HTTP server: %v", err)
		}
	}()

	// Wait for an interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	logrus.Info("Commencing graceful shutdown, press Ctrl+C again to force exiting")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := metricsSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Metrics server forced to shutdown: %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

}

// makeRouter initializes the application router and return it
func (s *server) makeRouter() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Health
	r.GET("/_/health", handlers.HealthHandler(s.state))
	r.GET("/.well-known/morty.json", handlers.WellknownHandler())

	// Functions
	r.GET("/functions", handlers.ListFunctionsHandler(s.state, s.orch))
	r.POST("/functions", handlers.CreateFunctionHandler(s.state, s.orch))
	r.Any("/functions/:name/:version/invoke", handlers.InvokeFunctionHandler(s.state, s.orch))

	return r
}

// getInitialState will try to retrieve the existing functions by calling the
// underlying orchestrator, and will populate the server state engine if
// functions are found.
func (s *server) getInitialState() {
	ctx := context.Background()
	if functions, err := s.orch.GetFunctions(ctx); err == nil {
		// If functions are found, we need to populate the state with them
		if errs := s.state.SetMultiple(ctx, functions); len(errs) > 0 {
			logrus.Warnf("Failed to populate state with existing functions: %v", err)
		}
	} else {
		logrus.Warnf("Failed to load existing functions, will start with an empty list: %v", err)
	}
}
