package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config config
type Config struct {
	Core     SectionCore
	Postgres SectionPostgres
	Redis    SectionRedis
}

// SectionCore is sub section of config.
type SectionCore struct {
	Mode   string
	Port   string
	Secret string
}

// SectionPostgres is sub section of env.
type SectionPostgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

// SectionRedis is sub section of config
type SectionRedis struct {
	Host       string
	Port       string
	Password   string
	DB         int
	MaxRetries int
}

// LoadConfig load env from file and read in environment variables that match
func LoadConfig() (Config, error) {
	var env Config

	// REQUIRED if the config file does not have the extension in the name
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// look for config in the working directory
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Core
	env.Core.Mode = viper.GetString("core_mode")
	env.Core.Port = viper.GetString("core_port")
	env.Core.Secret = viper.GetString("core_secret")

	// Postgres
	env.Postgres.Host = viper.GetString("postgres_host")
	env.Postgres.Port = viper.GetString("postgres_port")
	env.Postgres.User = viper.GetString("postgres_user")
	env.Postgres.Password = viper.GetString("postgres_password")
	env.Postgres.DB = viper.GetString("postgres_db")

	// redis
	env.Redis.Host = viper.GetString("redis_host")
	env.Redis.Port = viper.GetString("redis_port")
	env.Redis.DB = viper.GetInt("redis_db")
	env.Redis.Password = viper.GetString("redis_password")
	env.Redis.MaxRetries = viper.GetInt("redis_max_retries")

	return env, nil
}
