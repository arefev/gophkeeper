package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
	_ "github.com/joho/godotenv/autoload"
)

const (
	address       string = ":3200"
	logLevel      string = "info"
	tokenSecret   string = ""
	dbUser        string = ""
	dbPassword    string = ""
	dbName        string = ""
	dbHost        string = ""
	dbPort        string = ""
	tokenDuration int    = 60
)

type Config struct {
	Address       string `env:"ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
	TokenSecret   string `env:"TOKEN_SECRET"`
	DBUser        string `env:"DB_USER"`
	DBPassword    string `env:"DB_PASSWORD"`
	DBName        string `env:"DB_NAME"`
	DBHost        string `env:"DB_HOST"`
	DBPort        string `env:"DB_PORT"`
	TokenDuration int    `env:"TOKEN_DURATION"`
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
	f.StringVar(&cnf.DBUser, "db-user", dbUser, "user for db connection")
	f.StringVar(&cnf.DBPassword, "db-pwd", dbPassword, "password for db connection")
	f.StringVar(&cnf.DBName, "db-name", dbName, "name for db connection")
	f.StringVar(&cnf.DBHost, "db-host", dbHost, "host for db connection")
	f.StringVar(&cnf.DBPort, "db-port", dbPort, "port for db connection")
	f.StringVar(&cnf.TokenSecret, "t-secret", tokenSecret, "token secret")
	f.IntVar(&cnf.TokenDuration, "t-duration", tokenDuration, "token duration")
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
