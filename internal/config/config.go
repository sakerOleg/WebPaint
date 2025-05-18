package config

import "os"

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBName string
	Port   string
}

func Load() Config {
	return Config{
		DBUser: getEnv("DBUSER", ""),
		DBPass: getEnv("DBPASS", ""),
		DBHost: getEnv("DBHOST", ""),
		DBName: getEnv("DBNAME", ""),
		Port:   getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
