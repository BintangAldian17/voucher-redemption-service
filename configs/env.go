package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLICT_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8000"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBAddress:  getEnv("DB_ADDRESS", "127.0.0.1:3306"),
		DBName:     getEnv("DB_NAME", "otto_digital"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
