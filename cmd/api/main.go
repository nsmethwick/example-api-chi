package main

import (
	"context"
	"encoding/json"
	"errors"
	"example-api-chi/internal/config"
	"example-api-chi/internal/db"
	"example-api-chi/internal/repository"
	"example-api-chi/internal/routes"
	svc "example-api-chi/internal/service"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
)

func main() {
	flag.Parse()
	config.LoadEnvironmentVariables("")

	// NOTE: Moving the logic to `app()` allows us to honour `defer` calls.
	if err := app(); err != nil {
		os.Exit(1)
	}
}

func app() error {
	ctx := context.Background()
	apiConfig := config.SetUpAPI()
	dbConfig := config.SetUpDB()

	database, err := db.Init(dbConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			slog.WarnContext(ctx, "error closing database", "error", err)
		}
	}()

	slog.Info("Connected to MySQL database")

	repo, err := repository.NewRepository(database)
	if err != nil {
		log.Fatalf("Failed to create database repository: %v", err)
	}
	transactor, err := repository.NewTransactor(database)
	if err != nil {
		log.Fatalf("Failed to create transactor repository: %v", err)
	}
	service := svc.NewCustomerService(repo, transactor)

	apiAddr := apiConfig.Address

	slog.LogAttrs(ctx, slog.LevelInfo, "service_starting",
		slog.Int("gomaxprocs", runtime.GOMAXPROCS(0)),
		slog.String("api_server_addr", apiAddr),
	)

	g := new(errgroup.Group)
	timeout := time.Duration(10) * time.Second

	// Routing
	r := chi.NewRouter()

	// Route Registration

	// Catch-all Route
	r.Get("/", (func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "response"})
	}))

	// Register /customer
	r.Route("/customer", func(r chi.Router) {
		routes.CustomerRoutes(r, service)
	})

	// API Server
	apiServer := &http.Server{
		Addr:              apiAddr,
		Handler:           r,
		ReadHeaderTimeout: timeout, // avoid slowloris https://app.deepsource.com/directory/analyzers/go/issues/GO-S2114
	}

	g.Go(func() error {
		if err := apiServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.InfoContext(ctx, "api server stopped")
				return nil
			}
			err := fmt.Errorf("api server failed to listen and serve requests: %w", err)
			slog.LogAttrs(ctx, slog.LevelError, "error handling API requests",
				slog.Any("error", err),
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
		slog.LogAttrs(ctx, slog.LevelWarn, "http server graceful shutdown",
			slog.String("signal", signalType.String()),
		)
		// NOTE: Shutdown is given a separate context from the API server.
		if err := apiServer.Shutdown(context.Background()); err != nil {
			err := fmt.Errorf("unable to gracefully stop API server: %w", err)
			slog.LogAttrs(ctx, slog.LevelError, "error shutting down API server",
				slog.Any("error", err),
			)
			return err
		}
		return nil
	})

	return g.Wait()
}
