package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload"
)

const (
	address     string = ":3200"
	logLevel    string = "info"
	logFilePath string = ""
	tokenSecret string = ""
)

type Config struct {
	Address     string `env:"ADDRESS"`
	LogLevel    string `env:"LOG_LEVEL"`
	LogFilePath string `env:"LOG_FILE_PATH"`
	TokenSecret string `env:"TOKEN_SECRET"`
}

func NewConfig(params []string) (*Config, error) {
	cnf := &Config{}

	if err := cnf.initFlags(params); err != nil {
		return &Config{}, err
	}

	if err := cnf.initEnvs(); err != nil {
		return &Config{}, err
	}

	return cnf, nil
}

func (cnf *Config) initFlags(params []string) error {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", address, "address to run server")
	f.StringVar(&cnf.LogLevel, "l", logLevel, "log level")
	f.StringVar(&cnf.LogFilePath, "f", logLevel, "log file path")
	f.StringVar(&cnf.TokenSecret, "t-secret", tokenSecret, "token secret")
	if err := f.Parse(params); err != nil {
		return fmt.Errorf("InitFlags: parse flags fail: %w", err)
	}

	return nil
}

func (cnf *Config) initEnvs() error {
	if err := env.Parse(cnf); err != nil {
		return fmt.Errorf("InitEnvs: parse envs fail: %w", err)
	}

	return nil
}
