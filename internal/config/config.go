package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	LogFilePath    string
	MetricsPort    string
	DefaultBuckets []float64
}

func LoadConfig() *Config {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		log.Fatal("LOG_FILE_PATH environment variable is required")
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	buckets := []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0}
	if bucketStr := os.Getenv("METRICS_BUCKETS"); bucketStr != "" {
		buckets = parseBuckets(bucketStr)
	}

	return &Config{
		LogFilePath:    logFilePath,
		MetricsPort:    metricsPort,
		DefaultBuckets: buckets,
	}
}

func parseBuckets(bucketStr string) []float64 {
	var buckets []float64
	for _, b := range strings.Split(bucketStr, ",") {
		if v, err := strconv.ParseFloat(strings.TrimSpace(b), 64); err == nil {
			buckets = append(buckets, v)
		}
	}
	return buckets
}
