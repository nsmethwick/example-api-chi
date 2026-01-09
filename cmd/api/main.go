package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"

	h "github.com/erp-api/internal/http"
)

var db *sql.DB

func main() {
	flag.Parse()

	// NOTE: Moving the logic to `app()` allows us to honour `defer` calls.
	if err := app(); err != nil {
		os.Exit(1)
	}
}

func app() error {
	ctx := context.Background()
	// logging.Level.Set(slog.LevelWarn) // omit this line if you want slog.LevelInfo
	// logger := logging.NewLogger()

	// cfg, err := config.Load(*cfgPath)
	// if err != nil {
	// 	logger.LogAttrs(ctx, slog.LevelError, "config_load", slog.Any("error", fmt.Errorf("unable to load config: %w", err)))
	// }
	// if err := cfg.Validate(); err != nil {
	// 	logger.LogAttrs(ctx, slog.LevelError, "config_validate", slog.Any("error", fmt.Errorf("invalid config: %w", err)))
	// }

	g := new(errgroup.Group)
	timeout := time.Duration(10) * time.Second //nolint:gomnd

	// NOTE: All future log entries will include the image tag.
	// e.g. <DATE_TIME>-<GIT_BRANCH>-<GIT_COMMIT>
	// slog = slog.With(slog.String("build", ImageTag))

	// Routing
	mux := http.NewServeMux()

	// NOTE: We have separate pipelines.
	// One for the healthcheck and another for all other endpoints.
	// This is because we want to avoid logging/tracing the healthcheck endpoint.
	// For some apps, tracing the healthcheck can result in poor performance.

	pipelineHealth := middleware.New(middleware.AddRoutePattern(mux), metric.WithHealthCheckMetrics, h.WithHeaders(cfg))
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "image_tag": "%s"}`, ImageTag)
	})
	mux.Handle("/healthcheck", pipelineHealth.Decorate(healthHandler))

	// NOTE: The order of the middleware is deliberate.
	// AddRoutePattern updates the request context with the route pattern.
	// WithSpan starts a trace and updates the request context with a Trace ID.
	// WithRouteAndTrace attaches both the route pattern and Trace ID to the slog.
	// WithMetrics, WithRouteTag and WithParamAttrs all check for the route pattern added by WithRoutePattern.

	pipeline := middleware.New(
		middleware.AddRoutePattern(mux),
		traces.WithMessageEvents,
		logging.WithRouteAndTrace(logger),
		metric.WithMetrics,
		traces.WithRouteTag,
		traces.WithParamAttrs,
		h.WithHeaders(cfg),
	)

	exampleHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"endpoint": "example"}`)
	})
	mux.Handle("/example/{id}", pipeline.Decorate(exampleHandler))

	// API Server
	apiServer := &http.Server{
		Addr:              "localhost:3000",
		Handler:           mux,
		ReadHeaderTimeout: timeout, // avoid slowloris https://app.deepsource.com/directory/analyzers/go/issues/GO-S2114
	}
	g.Go(func() error {
		if err := apiServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.InfoContext(ctx, "api server stopped")
				return nil
			}
			err := fmt.Errorf("api server failed to listen and serve requests: %w", err)
			logger.LogAttrs(ctx, slog.LevelError, "error handling API requests",
				slog.Any("error", err),
				slog.String("metrics_addr", Addr),
			)
			return err
		}
		return nil
	})

	// Graceful shutdown
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		signalType := <-quit
		logger.LogAttrs(ctx, slog.LevelWarn, "http server graceful shutdown",
			slog.String("signal", signalType.String()),
		)
		// NOTE: Shutdown is given a separate context from the API server.
		if err := apiServer.Shutdown(context.Background()); err != nil { //nolint:contextcheck
			err := fmt.Errorf("unable to gracefully stop API server: %w", err)
			logger.LogAttrs(ctx, slog.LevelError, "error shutting down API server",
				slog.Any("error", err),
			)
			return err
		}
		if err := Server.Shutdown(context.Background()); err != nil { //nolint:contextcheck
			err := fmt.Errorf("unable to gracefully stop metrics server: %w", err)
			logger.LogAttrs(ctx, slog.LevelError, "error shutting down metrics server",
				slog.Any("error", err),
			)
			return err
		}
		return nil
	})

	return g.Wait() //nolint:wrapcheck
}

func run() error {
	logger := slog.Default()
	logger.Info("Starting API...")

	// do stuff here

	// db
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "root"
	cfg.Net = "tcp"
	cfg.Addr = "localhost:8888"
	cfg.DBName = "customers"

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// api setup
	// api := config.SetUpAPI()
	// address := api.Host + ":" + api.Port

	// r := chi.NewRouter()

	// r.Get("/", HandleRoot)

	// // grouped routes
	// r.Route("/customer", CustomerRoutes)

	// logger.Info("Listening at: " + address)
	// http.ListenAndServe(address, r)
	return nil
}

// routes for /customer
func CustomerRoutes(r chi.Router) {

	r.Use(middleware.Logger)

	r.Get("/", HandleListCustomers)           // GET /customer
	r.Get("/{customerID}", HandleGetCustomer) // GET /customer/123
	r.Post("/", HandlePostCustomer)           // POST /customer
}

func HandleListCustomers(w http.ResponseWriter, r *http.Request) {
	log.Printf("request path: %v", r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"data": "response body"}
	json.NewEncoder(w).Encode(response)
}

func HandleGetCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("request path: %v", r.URL)
	customerID := chi.URLParam(r, "customerID")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"customer": customerID}
	json.NewEncoder(w).Encode(response)
}

func HandlePostCustomer(w http.ResponseWriter, r *http.Request) {
	log.Printf("request path : %v", r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{"data": "POST response body"}
	json.NewEncoder(w).Encode(response)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("request path: %v", r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"data": "Welcome to our app!"}
	json.NewEncoder(w).Encode(response)
}
