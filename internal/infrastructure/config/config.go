package config

import (
	"os"
	"strconv"
	"time"
	"strings"
)

type Config struct {
	HTTPAddr string
	KafkaBrokers []string
	KafkaTopic string
	GeneratorEnabled bool
	GeneratorInterval time.Duration
	GeneratorUsers int
	GeneratorMovies int
}

func Load() *Config {
	return &Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		KafkaBrokers: splitEnv("KAFKA_BROKERS", ",", []string{"localhost:9092"}),
		KafkaTopic: getEnv("KAFKA_TOPIC", "movie-events"),
		GeneratorEnabled: getBoolEnv("GENERATOR_ENABLED", false),
		GeneratorInterval: getDurationEnv("GENERATOR_INTERVAL", "5s"),
		GeneratorUsers: getIntEnv("GENERATOR_USERS", 50),
		GeneratorMovies: getIntEnv("GENERATOR_MOVIES", 20),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return defaultValue
}

func splitEnv(key, sep string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var result []string
	for _, part := range strings.Split(value, sep) {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

func getDurationEnv(key, defaultValue string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		d, _ := time.ParseDuration(defaultValue)
		return d
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		d, _ := time.ParseDuration(defaultValue)
		return d
	}
	return d
}

func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}
