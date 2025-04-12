package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

const (
	address     string = "3200"
	databaseDSN string = ""
)

type Config struct {
	Address     string `env:"ADDRESS" json:"address"`
	DatabaseDSN string `env:"DATABASE_URI"`
}

func NewConfig(params []string) (Config, error) {
	cnf := Config{}

	if err := cnf.initFlags(params); err != nil {
		return Config{}, err
	}

	if err := cnf.initEnvs(); err != nil {
		return Config{}, err
	}

	return cnf, nil
}

func (cnf *Config) initFlags(params []string) error {
	f := flag.NewFlagSet("main", flag.ExitOnError)
	f.StringVar(&cnf.Address, "a", cnf.Address, "address to run server")
	f.StringVar(&cnf.DatabaseDSN, "d", databaseDSN, "db connection string")
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
