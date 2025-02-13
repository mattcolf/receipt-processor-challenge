package api

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment        string
	ServerHost         string
	ServerPort         int
	ServerBindAddress  string
	ServerWriteTimeout time.Duration
	ServerReadTimeout  time.Duration
}

func LoadConfig() *Config {
	host := GetEnvString("API_HOSTNAME", "0.0.0.0")
	port := GetEnvInt("API_PORT", 8080)

	return &Config{
		Environment:        GetEnvString("API_ENV", "production"),
		ServerHost:         host,
		ServerPort:         port,
		ServerBindAddress:  fmt.Sprintf("%s:%d", host, port),
		ServerWriteTimeout: GetEnvDuration("API_WRITE_TIMEOUT", 15*time.Second),
		ServerReadTimeout:  GetEnvDuration("API_READ_TIMEOUT", 15*time.Second),
	}
}

func GetEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		log.Printf("Missing environment variable %s, defaulting to %s instead.", key, defaultValue)
		return defaultValue
	}
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, error := strconv.Atoi(value); error == nil {
			return intValue
		} else {
			log.Printf("Invalid integer value for environment variable %s, defaulting to %d instead. %s", key, defaultValue, error)
			return defaultValue
		}
	} else {
		log.Printf("Missing environment variable %s, defaulting to %d instead.", key, defaultValue)
		return defaultValue
	}
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if durationValue, error := time.ParseDuration(value); error == nil {
			return durationValue
		} else {
			log.Printf("Invalid duration value for environment variable %s, defaulting to %s instead. %s", key, defaultValue, error)
			return defaultValue
		}
	} else {
		log.Printf("Missing environment variable %s, defaulting to %d instead.", key, defaultValue)
		return defaultValue
	}
}
