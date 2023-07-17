package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetConfig() *Config {
	loadEnv()

	return &Config{
		Port:             os.Getenv("PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresDb:       os.Getenv("POSTGRES_DB"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresSslMode:  os.Getenv("SSL_MODE"),
		RedisUrl:         os.Getenv("REDIS_URL"),
	}
}
