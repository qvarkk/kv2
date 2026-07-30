package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	LogLevel  EnvKey = "LOG_LEVEL"
	LogFormat EnvKey = "LOG_FORMAT"

	APIPort EnvKey = "API_PORT"
)

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

func LoadEnv() error {
	return godotenv.Load(".env")
}
