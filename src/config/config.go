package config

import "os"

type DatabaseInfo struct {
	HOST string
	PORT string
	NAME string
	USER string
	PASS string
}

type Configuration struct {
	Database DatabaseInfo
	Secret   string
}

var Env Configuration

func (cfg *Configuration) Init() {
	cfg.Database.HOST = getEnv("DB_HOST", "localhost")
	cfg.Database.PORT = getEnv("DB_PORT", "5432")
	cfg.Database.NAME = getEnv("DB_NAME", "")
	cfg.Database.USER = getEnv("DB_USER", "")
	cfg.Database.PASS = getEnv("DB_PASS", "")
	cfg.Secret = getEnv("SECRET", "mysecret")
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
