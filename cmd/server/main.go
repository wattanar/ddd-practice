package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"github.com/wattanar/taskmanager/api/spec"
	deptapp "github.com/wattanar/taskmanager/internal/application/department"
	taskapp "github.com/wattanar/taskmanager/internal/application/task"
	"github.com/wattanar/taskmanager/internal/infrastructure/http"
	"github.com/wattanar/taskmanager/internal/infrastructure/persistence/postgres"
)

func main() {
	fx.New(
		fx.Provide(
			loadConfig,
			newPool,
			postgres.NewTaskRepository,
			postgres.NewDepartmentRepository,
			taskapp.NewCreateTaskUseCase,
			taskapp.NewGetTaskUseCase,
			taskapp.NewListTasksUseCase,
			taskapp.NewUpdateTaskUseCase,
			taskapp.NewDeleteTaskUseCase,
			deptapp.NewCreateDepartmentUseCase,
			deptapp.NewGetDepartmentUseCase,
			deptapp.NewListDepartmentsUseCase,
			deptapp.NewUpdateDepartmentUseCase,
			deptapp.NewDeleteDepartmentUseCase,
			taskhttp.NewTaskHandler,
			taskhttp.NewDepartmentHandler,
			taskhttp.NewAPIHandler,
		),
		fx.Invoke(startServer),
	).Run()
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

func newPool(lc fx.Lifecycle, cfg config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func startServer(lc fx.Lifecycle, handler *taskhttp.APIHandler, cfg config) {
	mux := http.NewServeMux()
	specHandler := spec.HandlerFromMux(handler, mux)
	wrapped := taskhttp.RecoveryMiddleware(taskhttp.Middleware(specHandler))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrapped,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.Info("server starting", "port", cfg.Port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("shutting down server...")
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
