package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl          string
	JWTSecret      string
	Port           string
	GoogleClientID string
	AppEnv         string
	TokenTTL       int
}

func Load() *Config {
	_ = godotenv.Load() // ignore error: env may be set in .env
	cfg := &Config{
		DBUrl:          getenv("DATABASE_URL", "postgres://postgres:25244002myKaffa@localhost:5432/ops_system?sslmode=disable"),
		JWTSecret:      getenv("JWT_SECRET", "bdd158a2de8e3401e062cc0e7523384f"),
		Port:           getenv("PORT", "8080"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		AppEnv:         getenv("APP_ENV", "Development"),
		TokenTTL:       getint("TOKEN_TTL_HOURS", 24),
	}
	return cfg
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func getint(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return d

}
