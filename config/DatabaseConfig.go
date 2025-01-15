package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	MongoURI string
	Port     string
}

func LoadConfig() *Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	// Ensure the URI starts with valid schemes
	if !(strings.HasPrefix(mongoURI, "mongodb://") || strings.HasPrefix(mongoURI, "mongodb+srv://")) {
		log.Fatal("Invalid MongoDB URI, must start with 'mongodb://' or 'mongodb+srv://'")
	}
	return &Config{
		MongoURI: mongoURI,
		Port:     os.Getenv("PORT"),
	}
}
