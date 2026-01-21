package config

import "os"

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "admin"
	defaultDBPassword = "glc_eyJvIjoiMTYzMTM5MSIsIm4iOiJzdGFjay0xNDg0MzE0LWFsbG95LWluc3RhZ3JhbSIsImsiOiI1NVMyMEdFRzZOM3JKRzBxWTY0TUM0bGUiLCJtIjp7InIiOiJwcm9kLWFwLXNvdXRoLTEifX0="
	defaultDBName     = "grademanagement"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

var (
	dbhost     = getEnv("DB_HOST", defaultDBHost)
	dbport     = getEnv("DB_PORT", defaultDBPort)
	dbuser     = getEnv("DB_USER", defaultDBUser)
	dbpassword = getEnv("DB_PASSWORD", defaultDBPassword)
	dbname     = getEnv("DB_NAME", defaultDBName)
	dbpassword = os.Getenv("DB_PASSWORD")
	dbname     = getenvOrDefault("DB_NAME", "grademanagement")
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
