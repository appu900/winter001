package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env          string
	HttpAddr     string
	DatabaseURL  string
	RedisURL     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdealTimeOut time.Duration
}

func Load() (*Config, error) {
	var errs []error
	cfg := &Config{
		Env:      GetEnv("APP_ENV", "development"),
		HttpAddr: GetEnv("HTTP_PORT", ":8080"),
	}
	cfg.DatabaseURL = mustEnv("DATABASE_URL", &errs)
	cfg.RedisURL = mustEnv("REDIS_URL", &errs)

	cfg.ReadTimeout = getDuration("HTTP_READ_TIMEOUT", 5*time.Second, &errs)
	cfg.WriteTimeout = getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second, &errs)
	cfg.IdealTimeOut = getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second, &errs)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil

}

func GetEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string, errs *[]error) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		*errs = append(*errs, fmt.Errorf("%s is required", key))
	}
	return v
}

func getDuration(key string, fallbackDuration time.Duration, errs *[]error) time.Duration {
	value, ok := os.LookupEnv(key)
	if value == "" || !ok {
		return fallbackDuration
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("invalid duration"))
		return fallbackDuration
	}
	return d
}
func getInt(key string, fallback int, errs *[]error) int {
	value, ok := os.LookupEnv(key)
	if value == "" || !ok {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: invalid int format", key))
		log.Println("falling back to fallback valuek")
		return fallback
	}
	return n
}
