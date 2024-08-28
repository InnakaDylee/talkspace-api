package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	POSTGRESQL    PostgreSQLConfig
	REDIS         RedisConfig
	ELASTICSEARCH ElasticsearchConfig
	CLOUDSTORAGE  CloudStorageConfig
	MIDTRANS      MidtransConfig
	SMTP          SMTPConfig
	OPENAI        OpenAIConfig
	JWT           JWTConfig
	SERVER        ServerConfig
}

type (
	PostgreSQLConfig struct {
		POSTGRESQL_USER string
		POSTGRESQL_PASS string
		POSTGRESQL_HOST string
		POSTGRESQL_PORT string
		POSTGRESQL_NAME string
	}

	RedisConfig struct {
		REDIS_HOST string
		REDIS_PORT string
		REDIS_DB   string
		REDIS_PASS string
	}

	ElasticsearchConfig struct {
		ELASTICSEARCH_URL  string
		ELASTICSEARCH_USER string
		ELASTICSEARCH_PASS string
	}

	CloudStorageConfig struct {
		AWS_ACCESS_KEY_ID     string
		AWS_SECRET_ACCESS_KEY string
		AWS_REGION            string
		AWS_BUCKET_NAME       string
	}

	MidtransConfig struct {
		MIDTRANS_SERVER_KEY string
		MIDTRANS_CLIENT_KEY string
	}

	OpenAIConfig struct {
		OPENAI_API_KEY string
	}

	SMTPConfig struct {
		SMTP_USER string
		SMTP_PASS string
		SMTP_PORT string
		SMTP_HOST string
	}

	ServerConfig struct {
		SERVER_HOST string
		SERVER_PORT string
	}

	JWTConfig struct {
		JWT_SECRET string
	}
)

func LoadConfig() (*Configuration, error) {

	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load environment variables from .env file: %w", err)
		}
	} else {
		fmt.Println("warning: .env file not found. make sure environment variables are set")
	}

	return &Configuration{
		POSTGRESQL: PostgreSQLConfig{
			POSTGRESQL_USER: os.Getenv("POSTGRESQL_USER"),
			POSTGRESQL_PASS: os.Getenv("POSTGRESQL_PASS"),
			POSTGRESQL_HOST: os.Getenv("POSTGRESQL_HOST"),
			POSTGRESQL_PORT: os.Getenv("POSTGRESQL_PORT"),
			POSTGRESQL_NAME: os.Getenv("POSTGRESQL_NAME"),
		},
		REDIS: RedisConfig{
			REDIS_HOST: os.Getenv("REDIS_HOST"),
			REDIS_PORT: os.Getenv("REDIS_PORT"),
			REDIS_DB:   os.Getenv("REDIS_DB"),
			REDIS_PASS: os.Getenv("REDIS_PASS"),
		},
		ELASTICSEARCH: ElasticsearchConfig{
			ELASTICSEARCH_URL:  os.Getenv("ELASTICSEARCH_URL"),
			ELASTICSEARCH_USER: os.Getenv("ELASTICSEARCH_USER"),
			ELASTICSEARCH_PASS: os.Getenv("ELASTICSEARCH_PASS"),
		},
		CLOUDSTORAGE: CloudStorageConfig{
			AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AWS_REGION:            os.Getenv("AWS_REGION"),
			AWS_BUCKET_NAME:       os.Getenv("AWS_BUCKET_NAME"),
		},
		MIDTRANS: MidtransConfig{
			MIDTRANS_SERVER_KEY: os.Getenv("MIDTRANS_SERVER_KEY"),
			MIDTRANS_CLIENT_KEY: os.Getenv("MIDTRANS_CLIENT_KEY"),
		},
		SMTP: SMTPConfig{
			SMTP_USER: os.Getenv("SMTP_USER"),
			SMTP_PASS: os.Getenv("SMTP_PASS"),
			SMTP_PORT: os.Getenv("SMTP_PORT"),
			SMTP_HOST: os.Getenv("SMTP_HOST"),
		},
		OPENAI: OpenAIConfig{
			OPENAI_API_KEY: os.Getenv("OPENAI_API_KEY"),
		},
		SERVER: ServerConfig{
			SERVER_HOST: os.Getenv("SERVER_HOST"),
			SERVER_PORT: os.Getenv("SERVER_PORT"),
		},
		JWT: JWTConfig{
			JWT_SECRET: os.Getenv("JWT_SECRET"),
		},
	}, nil
}
