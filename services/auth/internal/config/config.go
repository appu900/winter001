package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type ApplicationConfig struct {
	Env                string
	ServerAddrs        string
	DatabaseUrl        string
	RedisAddress       string // host:port, built from REDIS_HOST + REDIS_PORT
	RedisUsername      string
	RedisPassword      string
	MaxPoolsize        int
	MaxIdealConnection int
	MinIdealConnection int
	MaxPoolTimeout     time.Duration
	WriteTimeout       time.Duration
	ReadTimeout        time.Duration
}

func Load() (*ApplicationConfig, error) {
	var errs []error
	appcfg := &ApplicationConfig{
		Env:         GetEnv("APP_ENV", "development"),
		ServerAddrs: GetEnv("AUTH_HTTP_ADDR", ":8081"),
	}

	// database
	appcfg.DatabaseUrl = mustEnv("AUTH_DATABASE_URL", &errs)

	// redis
	redisHost := GetEnv("REDIS_HOST", "localhost")
	redisPort := GetEnv("REDIS_PORT", "6379")
	appcfg.RedisAddress = net.JoinHostPort(redisHost, redisPort)
	appcfg.RedisUsername = GetEnv("REDIS_USERNAME", "") // optional: only for Redis ACL users
	appcfg.RedisPassword = GetEnv("REDIS_PASSWORD", "") // optional: empty for local Redis
	appcfg.MaxPoolsize = getInt("REDIS_POOL_SIZE", 20, &errs)
	appcfg.MaxIdealConnection = getInt("REDIS_MAX_IDLE_CONNS", 10, &errs)
	appcfg.MinIdealConnection = getInt("REDIS_MIN_IDLE_CONNS", 5, &errs)
	appcfg.MaxPoolTimeout = getDuration("REDIS_POOL_TIMEOUT", 5*time.Second, &errs)

	// http server
	appcfg.ReadTimeout = getDuration("HTTP_READ_TIMEOUT", 5*time.Second, &errs)
	appcfg.WriteTimeout = getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second, &errs)

	if appcfg.MinIdealConnection > appcfg.MaxIdealConnection {
		errs = append(errs, fmt.Errorf("REDIS_MIN_IDLE_CONNS (%d) must be <= REDIS_MAX_IDLE_CONNS (%d)",
			appcfg.MinIdealConnection, appcfg.MaxIdealConnection))
	}
	if appcfg.MaxIdealConnection > appcfg.MaxPoolsize {
		errs = append(errs, fmt.Errorf("REDIS_MAX_IDLE_CONNS (%d) must be <= REDIS_POOL_SIZE (%d)",
			appcfg.MaxIdealConnection, appcfg.MaxPoolsize))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return appcfg, nil
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
		*errs = append(*errs, fmt.Errorf("%s: invalid duration %q", key, value))
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
		*errs = append(*errs, fmt.Errorf("%s: invalid int %q", key, value))
		return fallback
	}
	return n
}
