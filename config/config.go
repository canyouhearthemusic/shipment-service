package config

import "os"

type (
	GRPC struct {
		Port string
	}

	Config struct {
		GRPC GRPC
	}
)

func Load() Config {
	return Config{
		GRPC: GRPC{
			Port: getEnv("GRPC_PORT", "50051"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
