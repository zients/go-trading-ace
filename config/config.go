package config

import (
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Infura   InfuraConfig   `mapstructure:"infura"`
}

type ServerConfig struct {
	Port                  int `mapstructure:"port"`
	RequestTimeoutSeconds int `mapstructure:"request_timeout_seconds"`
}

func (s ServerConfig) RequestTimeout() time.Duration {
	if s.RequestTimeoutSeconds <= 0 {
		return 10 * time.Second
	}

	return time.Duration(s.RequestTimeoutSeconds) * time.Second
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Prefix string `mapstructure:"prefix"`
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
}

type InfuraConfig struct {
	Key string `mapstructure:"key"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.SetConfigType("yaml")

	return loadConfig(v)
}

func loadConfig(v *viper.Viper) (*Config, error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.Set("server.request_timeout_seconds", parseRequestTimeoutSeconds(v.Get("server.request_timeout_seconds")))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseRequestTimeoutSeconds(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}

		return parsed
	default:
		return 0
	}
}
