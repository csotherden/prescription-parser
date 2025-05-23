package main

import (
	"context"
	"flag"
	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/server"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var envFile string

	flag.StringVar(&envFile, "env", ".env", "env file path")
	flag.Parse()

	// Load environment variables from .env file
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create server config
	cfg := config.NewConfig()

	// Create and configure server
	srv, err := server.NewServer(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c
	logger.Info("Received signal, shutting down server", zap.String("signal", sig.String()))

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server gracefully stopped")
}
