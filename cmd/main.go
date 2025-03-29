package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fakubwoy/go-file-share/api"
	"github.com/fakubwoy/go-file-share/internal/config"
	"github.com/fakubwoy/go-file-share/internal/database"
	"github.com/fakubwoy/go-file-share/internal/storage"
	"github.com/fakubwoy/go-file-share/internal/worker"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	var fileStorage storage.Storage
	if cfg.S3Enabled {
		fileStorage, err = storage.NewS3Storage(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize S3 storage: %v", err)
		}
	} else {
		fileStorage, err = storage.NewLocalStorage(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize local storage: %v", err)
		}
	}

	cleanupWorker := worker.NewCleanupWorker(db, fileStorage, 1*time.Hour)
	go cleanupWorker.Start()

	router := api.SetupRoutes(db, rdb, cfg, fileStorage)
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Printf("Server started on port %s", cfg.ServerPort)

	<-done
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited properly")
}
