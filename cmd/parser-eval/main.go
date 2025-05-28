package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/csotherden/prescription-parser/pkg/config"
	"github.com/csotherden/prescription-parser/pkg/datastore"
	"github.com/csotherden/prescription-parser/pkg/jobs"
	"github.com/csotherden/prescription-parser/pkg/models"
	"github.com/csotherden/prescription-parser/pkg/parser"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	var envFile string
	var testPdf string
	var testJson string
	var iterations int

	flag.StringVar(&envFile, "env", ".env", "env file path")
	flag.StringVar(&testPdf, "pdf", "", "test PDF file path")
	flag.StringVar(&testJson, "json", "", "test JSON file path")
	flag.IntVar(&iterations, "iterations", 1, "number of iterations to run")
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

	// Initialize datastore
	ds, err := datastore.NewPgEntDatastore(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize datastore", zap.Error(err))
	}

	// Initialize parser with the appropriate backend
	parserInstance, err := parser.NewParser(cfg, ds, logger)
	if err != nil {
		logger.Fatal("Failed to initialize parser", zap.Error(err))
	}

	fileName := filepath.Base(testPdf)

	// Read the test PDF file
	inputPdf, err := os.ReadFile(testPdf)
	if err != nil {
		logger.Fatal("Failed to read test PDF file", zap.Error(err))
	}

	// Read the test JSON file
	expectedJson, err := os.ReadFile(testJson)
	if err != nil {
		logger.Fatal("Failed to read test JSON file", zap.Error(err))
	}

	wg := sync.WaitGroup{}

	wg.Add(iterations)

	jobIds := []string{}

	// Run the parser
	for i := 0; i < iterations; i++ {
		jobId, err := parserInstance.ParseImage(context.Background(), fileName, bytes.NewReader(inputPdf))
		if err != nil {
			logger.Fatal("Failed to parse PDF", zap.Error(err))
		}

		jobIds = append(jobIds, jobId)

		go waitForResult(&wg, logger, jobId)
	}

	wg.Wait()

	for _, jobId := range jobIds {
		jobStatus, ok := jobs.GlobalTracker.GetJob(jobId)
		if !ok {
			logger.Error("Job not found", zap.String("jobId", jobId))
			continue
		}

		if jobStatus.Status == jobs.JobStatusFailed {
			logger.Error("Job failed", zap.String("jobId", jobId), zap.String("error", jobStatus.Error))
			continue
		}

		rx, ok := jobStatus.Result.(models.Prescription)
		if !ok {
			logger.Error("Job result is not a prescription", zap.String("jobId", jobId))
			continue
		}

		rxJson, err := json.Marshal(rx)
		if err != nil {
			logger.Error("Failed to marshal prescription", zap.String("jobId", jobId), zap.Error(err))
			continue
		}

		score, err := parserInstance.ScoreResult(context.Background(), string(expectedJson), string(rxJson))
		if err != nil {
			logger.Error("Failed to score result", zap.String("jobId", jobId), zap.Error(err))
			continue
		}

		fmt.Printf("Filename: %s - Job ID: %s\nScore: %.2f%% - (%.2f / %d)\nFeedback:\n%s\n\n", fileName, jobId, score.OverallScorePercentage, score.TotalAwardedPoints, int(score.TotalPossiblePoints), score.SummaryCritique)
	}
}

func waitForResult(wg *sync.WaitGroup, logger *zap.Logger, jobId string) {
	deadline := time.Now().Add(2 * time.Minute)

	for {
		time.Sleep(5 * time.Second)

		status, ok := jobs.GlobalTracker.GetJob(jobId)
		if !ok {
			logger.Fatal("Job not found", zap.String("jobId", jobId))
			break
		}

		if status.Status != jobs.JobStatusPending && status.Status != jobs.JobStatusProcessing {
			break
		}

		// Check if deadline exceeded
		if time.Now().After(deadline) {
			logger.Error("Job processing timeout exceeded (2 minutes)", zap.String("jobId", jobId))
			break
		}
	}

	wg.Done()
}
