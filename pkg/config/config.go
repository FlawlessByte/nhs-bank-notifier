// pkg/config/config.go
package config

import (
	"nhs-bank-notifier/pkg/logger"
	"os"
	"strconv"
	"time"
)

type Config struct {
	NHSUsername       string
	NHSPassword       string
	LoginURL          string
	TelegramToken     string
	TelegramChatID    string
	MaxTTL            time.Duration
	CheckIntervalMins int
	LogLevel          string
	NHSUnit           string
}

func LoadConfig() *Config {
	log := logger.GetLogger()
	// Read environment variables and set defaults if needed
	maxTTL, err := time.ParseDuration(getEnv("MAX_TTL", "336h")) //Default 2 weeks
	if err != nil {
		log.Fatalf("Invalid MAX_TTL: %v", err)
	}

	checkIntervalMins, err := strconv.Atoi(getEnv("CHECK_INTERVAL_MINS", "10"))
	if err != nil {
		log.Fatalf("Invalid CHECK_INTERVAL_MINS: %v", err)
	}

	return &Config{
		NHSUsername:       getEnv("NHS_USERNAME", ""),
		NHSPassword:       getEnv("NHS_PASSWORD", ""),
		LoginURL:          getEnv("LOGIN_URL", "https://ich.allocate-cloud.co.uk/EmployeeOnlineHealth/ICHLIVE/Login"),
		TelegramToken:     getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramChatID:    getEnv("TELEGRAM_CHAT_ID", ""),
		MaxTTL:            maxTTL,
		CheckIntervalMins: checkIntervalMins,
		LogLevel:          getEnv("LOG_LEVEL", "WARN"),
		NHSUnit:           getEnv("NHS_UNIT", "Intensive Care"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
