package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string       `yaml:"env" env:"ENV" env-default:"local"`
	GRPC   GRPCConfig   `yaml:"grpc"`
	DB     DBConfig     `yaml:"postgres"`
	Redis  RedisConfig  `yaml:"redis"`
	Rabbit RabbitConfig `yaml:"rabbitmq"`
	S3     S3Config     `yaml:"s3"`
}

// GRPC сервер
type GRPCConfig struct {
	Host    string        `yaml:"host" env:"GRPC_HOST" env-default:"localhost"`
	Port    int           `yaml:"port" env:"GRPC_PORT" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"10s"`
}

// PostgreSQL
type DBConfig struct {
	Host        string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port        int           `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name        string        `yaml:"name" env:"DB_NAME" env-default:"tournaments"`
	User        string        `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password    string        `yaml:"password" env:"DB_PASSWORD" env-default:""`
	SSLMode     string        `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
	MaxConns    int           `yaml:"max_conns" env:"DB_MAX_CONNS" env-default:"10"`
	MinConns    int           `yaml:"min_conns" env:"DB_MIN_CONNS" env-default:"2"`
	MaxConnIdle time.Duration `yaml:"max_conn_idle" env:"DB_MAX_CONN_IDLE" env-default:"30m"`
	ConnTimeout time.Duration `yaml:"conn_timeout" env:"DB_CONN_TIMEOUT" env-default:"5s"`
}

// Redis
type RedisConfig struct {
}

// RabbitMQ
type RabbitConfig struct {
}

// S3/MinIO
type S3Config struct {
}

func MustRead() *Config {
	path := fetchConfigPath()
	var cfg Config

	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(fmt.Sprintf("config file does not exist: %s", path))
		}

		// cleanenv.ReadConfig читает сначала файл, потом env-переменные
		err := cleanenv.ReadConfig(path, &cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to load config: %v", err))
		}
	} else {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to read environment variables: %v", err))
		}
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config")
	flag.Parse()

	if path == "" {
		return os.Getenv("CONFIG_PATH")
	}

	return path
}
