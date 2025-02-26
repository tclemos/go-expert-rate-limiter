package config

import (
	"time"

	"github.com/spf13/viper"
)

type WebServerConfig struct {
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

	MaxRequestsPerSecondPerIP     uint64        `mapstructure:"MAX_REQUESTS_PER_SECOND_PER_IP"`
	MaxRequestsPerSecondPerAPIKey uint64        `mapstructure:"MAX_REQUESTS_PER_SECOND_PER_API_TOKEN"`
	BanDuration                   time.Duration `mapstructure:"BAN_DURATION"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort int    `mapstructure:"REDIS_PORT"`
	RedisDB   int    `mapstructure:"REDIS_DB"`
	RedisPass string `mapstructure:"REDIS_PASS"`
}

func LoadConfig[T any](path string, config *T) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
