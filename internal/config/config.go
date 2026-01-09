package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type (
	// Environment represents the cluster environment.
	Environment string
)

var (
	// EnvironmentLocal is the local developer environment.
	EnvironmentLocal Environment = "local"
	// EnvironmentDev is the dev cluster environment.
	EnvironmentDev Environment = "dev"
	// EnvironmentStaging is the stage cluster environment.
	EnvironmentStaging Environment = "stg"
	// EnvironmentProduction is the production cluster environment.
	EnvironmentProduction Environment = "prd"
	// AllEnvironments is an aggregate of all environments.
	AllEnvironments = []Environment{EnvironmentLocal, EnvironmentDev, EnvironmentStaging, EnvironmentProduction}
)

// GetEnvironment returns the environment object.
func GetEnvironment(v string) Environment {
	for _, e := range AllEnvironments {
		if e.String() == v {
			return e
		}
	}
	return EnvironmentDev
}

// String returns the name of the environment
func (e Environment) String() string {
	return string(e)
}

// Config contains all the values necessary to construct a Server
type Config struct {
	Addr            string
	CacheDefaultTTL int
	Environment     Environment // Environment name, e.g. "dev"
	Metrics         Metrics
	OtelDisabled    bool
	Secrets         Secrets
}

// Metrics contains the metric server address and port.
type Metrics struct {
	Host string
	Port int
}

// Secrets contains Vault secrets exposed to the config.
//
// NOTE: The `Foo` secret here is used as an example you would modify/remove.
// Refer to ../../docs/secrets.md for details.
type Secrets struct {
	Foo string
}

// Load reads json from cfgPath and returns a populated Config
func Load(cfgPath string) (*Config, error) {
	var cfg Config
	file, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open cfg file '%s': %v", cfgPath, err)
	}
	if err = json.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal cfg file '%s': %v", cfgPath, err)
	}
	return &cfg, nil
}

// Validate checks specific fields contain valid values.
func (cfg *Config) Validate() error {
	var result []error
	if cfg.Addr == "" {
		result = append(result, errors.New("missing Addr config"))
	}
	if cfg.Environment == "" {
		result = append(result, errors.New("missing Environment config"))
	}
	validEnvironment := false
	for _, v := range AllEnvironments {
		if v == cfg.Environment {
			validEnvironment = true
		}
	}
	if !validEnvironment {
		result = append(result, fmt.Errorf("invalid Environment config: %q", cfg.Environment))
	}
	return errors.Join(result...)
}
