package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wattanar/taskmanager/api/spec"
	taskapp "github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/infrastructure/http"
	"github.com/wattanar/taskmanager/internal/infrastructure/persistence/postgres"
)

func main() {
	cfg := loadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	repo := postgres.NewTaskRepository(pool)

	createUC := taskapp.NewCreateTaskUseCase(repo)
	getUC := taskapp.NewGetTaskUseCase(repo)
	listUC := taskapp.NewListTasksUseCase(repo)
	updateUC := taskapp.NewUpdateTaskUseCase(repo)
	deleteUC := taskapp.NewDeleteTaskUseCase(repo)

	handler := taskhttp.NewTaskHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mux := http.NewServeMux()
	specHandler := spec.HandlerFromMux(handler, mux)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      taskhttp.RecoveryMiddleware(taskhttp.Middleware(specHandler)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}

type config struct {
	DatabaseURL string
	Port        string
}

func loadConfig() config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://taskmanager:taskmanager@localhost:5432/taskmanager?sslmode=disable"
	}

	return config{
		DatabaseURL: dbURL,
		Port:        port,
	}
}
