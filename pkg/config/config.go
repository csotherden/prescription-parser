// Package config provides configuration management for the prescription parser application.
// It handles loading environment variables and providing a central configuration structure.
package config

import (
	"os"
	"time"
)

// Config holds all application configuration settings.
// Values are loaded from environment variables when the application starts.
type Config struct {
	Host             string        // Host address for the HTTP server
	Port             string        // Port for the HTTP server
	ReadTimeout      time.Duration // Maximum duration for reading the entire request
	WriteTimeout     time.Duration // Maximum duration for writing the response
	IdleTimeout      time.Duration // Maximum duration to wait for the next request
	DatabaseHost     string        // Host address of the database server
	DatabasePort     string        // Port of the database server
	DatabaseName     string        // Name of the database to connect to
	DatabaseUser     string        // Username for database authentication
	DatabasePassword string        // Password for database authentication
	RunMigrations    bool          // Whether to run database migrations on startup
	OpenAIAPIKey     string        // API key for OpenAI services
	GeminiAPIKey     string        // API key for Gemini services
	ParserBackend    string        // Backend to use for prescription parsing ("OpenAI" or "Gemini")
}

// NewConfig creates a new Config instance with values loaded from environment variables.
// Default values are provided for some fields when environment variables are not set.
// If PARSER_BACKEND is not set, it will default to "OpenAI" or "Gemini" based on which API key is available.
func NewConfig() Config {
	// Load database configuration from environment
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Load AI service API keys
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	// Determine which parser backend to use
	parserBackend := os.Getenv("PARSER_BACKEND")

	// Auto-select parser backend based on available API keys if not explicitly set
	if parserBackend == "" && openAIAPIKey != "" {
		parserBackend = "OpenAI"
	} else if parserBackend == "" && geminiAPIKey != "" {
		parserBackend = "Gemini"
	}

	// Set default server port if not specified
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Set default server host if not specified
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	// Return the populated configuration
	return Config{
		Host:             host,
		Port:             port,
		ReadTimeout:      15 * time.Second,
		WriteTimeout:     15 * time.Second,
		IdleTimeout:      60 * time.Second,
		DatabaseHost:     dbHost,
		DatabasePort:     dbPort,
		DatabaseName:     dbName,
		DatabaseUser:     dbUser,
		DatabasePassword: dbPass,
		RunMigrations:    false,
		OpenAIAPIKey:     openAIAPIKey,
		GeminiAPIKey:     geminiAPIKey,
		ParserBackend:    parserBackend,
	}
}
