package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/soumyadri/todo-backend/internal/config"
	"github.com/soumyadri/todo-backend/internal/storage/sqlite"
	todo "github.com/soumyadri/todo-backend/internal/http/handlers/todo"
)

func main() {
	// Configuration Loading
	cfg := config.MustLoad()

	// Database Connection
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connection established", slog.String("env", cfg.Env), slog.String("storage_path", cfg.StoragePath))
	defer storage.Db.Close()

	// Setup Router
	http.NewServeMux()

	// Controller
	http.HandleFunc("POST /api/create/todo", todo.NewTodos(storage))
	http.HandleFunc("GET /api/todo", todo.GetTodos(storage))
	http.HandleFunc("GET /api/todo/{id}", todo.GetTodoById(storage))
	http.HandleFunc("PUT /api/todo/{id}", todo.UpdateTodo(storage))

	slog.Info("Server started", "address", cfg.HTTPServer.Address)

	// Setup gracefull shut down
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM);

	// Run server
	go func ()  {
		// start server
		if err := http.ListenAndServe(cfg.HTTPServer.Address, nil); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		} else {
			fmt.Printf("Server is running on %s\n", cfg.HTTPServer.Address)
		}
	}();

	// Wait for shutdown signal
	<-done
	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: http.DefaultServeMux,
	}
	error := server.Shutdown(ctx)

	// Check for error during shutdown
	if error != nil {
		slog.Error("Failed to shutdown server", "error", error)
	} else {
		slog.Info("Server shutdown successfully")
	}
}