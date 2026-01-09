package mcp

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// LoadConfigFromEnv loads MCP configuration from environment variables
func LoadConfigFromEnv() (Config, error) {
	baseURL := os.Getenv("MCP_BASE_URL")
	if baseURL == "" {
		return Config{}, fmt.Errorf("MCP_BASE_URL environment variable is required")
	}

	apiKey := os.Getenv("MCP_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("MCP_API_KEY environment variable is required")
	}

	config := Config{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}

	// Optional timeout configuration
	if timeoutStr := os.Getenv("MCP_TIMEOUT"); timeoutStr != "" {
		timeoutSec, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return Config{}, fmt.Errorf("invalid MCP_TIMEOUT value: %w", err)
		}
		if timeoutSec <= 0 {
			return Config{}, fmt.Errorf("MCP_TIMEOUT must be positive, got: %d", timeoutSec)
		}
		config.Timeout = time.Duration(timeoutSec) * time.Second
	}

	// Optional max retries configuration
	if retriesStr := os.Getenv("MCP_MAX_RETRIES"); retriesStr != "" {
		retries, err := strconv.Atoi(retriesStr)
		if err != nil {
			return Config{}, fmt.Errorf("invalid MCP_MAX_RETRIES value: %w", err)
		}
		if retries < 0 {
			return Config{}, fmt.Errorf("MCP_MAX_RETRIES cannot be negative, got: %d", retries)
		}
		config.MaxRetries = retries
	}

	// Optional logging configuration
	if loggingStr := os.Getenv("MCP_ENABLE_LOGGING"); loggingStr != "" {
		logging, err := strconv.ParseBool(loggingStr)
		if err != nil {
			return Config{}, fmt.Errorf("invalid MCP_ENABLE_LOGGING value: %w", err)
		}
		config.EnableLogging = logging
	}

	return config, nil
}
