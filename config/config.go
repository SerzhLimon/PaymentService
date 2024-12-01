package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	filepath = "config/config.env"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Password string `json:"password"`
}

type Config struct {
	Postgres PostgresConfig `json:"postgres"`
}

func LoadConfig() Config {
	if err := godotenv.Load(filepath); err != nil {
		log.Printf("No .env file found, falling back to environment variables: %v", err)
	}

	var config Config
	config.Postgres = PostgresConfig{
		Host:     getEnv("POSTGRES_HOST"),
		Port:     getEnv("POSTGRES_PORT"),
		User:     getEnv("POSTGRES_USER"),
		DBName:   getEnv("POSTGRES_DBNAME"),
		SSLMode:  getEnv("POSTGRES_SSLMODE"),
		Password: getEnv("POSTGRES_PASSWORD"),
	}

	return config
}

func getEnv(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}