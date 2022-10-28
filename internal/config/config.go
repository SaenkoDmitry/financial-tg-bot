package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const configFile = "data/config.yaml"

type Config struct {
	Token                       string        `yaml:"token"`
	AbstractAPIKey              string        `yaml:"abstract_api_key"`
	RatesCacheDefaultExpiration time.Duration `yaml:"rates_cache_default_expiration"`
	RatesCacheCleanupInterval   time.Duration `yaml:"rates_cache_cleanup_interval"`
	PostgresUser                string        `yaml:"postgres_user"`
	PostgresPassword            string        `yaml:"postgres_password"`
	PostgresDB                  string        `yaml:"postgres_db"`
	PostgresHost                string        `yaml:"postgres_host"`
	PostgresPort                string        `yaml:"postgres_port"`
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) AbstractAPIKey() string {
	return s.config.AbstractAPIKey
}

func (s *Service) RatesCacheDefaultExpiration() time.Duration {
	return s.config.RatesCacheDefaultExpiration
}

func (s *Service) RatesCacheCleanupInterval() time.Duration {
	return s.config.RatesCacheCleanupInterval
}

func (s *Service) PostgresUser() string {
	return s.config.PostgresUser
}

func (s *Service) PostgresPassword() string {
	return s.config.PostgresPassword
}

func (s *Service) PostgresDB() string {
	return s.config.PostgresDB
}

func (s *Service) PostgresHost() string {
	return s.config.PostgresHost
}

func (s *Service) PostgresPort() string {
	return s.config.PostgresPort
}
