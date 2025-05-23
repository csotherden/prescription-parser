package config

import (
	"os"
	"time"
)

type Config struct {
	Host             string
	Port             string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	RunMigrations    bool
	OpenAIAPIKey     string
	GeminiAPIKey     string
	ParserBackend    string
}

func NewConfig() Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	parserBackend := os.Getenv("PARSER_BACKEND")

	if parserBackend == "" && openAIAPIKey != "" {
		parserBackend = "OpenAI"
	} else if parserBackend == "" && geminiAPIKey != "" {
		parserBackend = "Gemini"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

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
