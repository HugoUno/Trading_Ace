package config

import (
	"trading-ace/pkg/database"
	"trading-ace/pkg/ethereum"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database database.Config
	Ethereum ethereum.Config
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type EthereumConfig struct {
	NodeURL     string
	PoolAddress string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
