// cmd/nhs-bank-notifier/main.go
package main

import (
	"time"

	"nhs-bank-notifier/pkg/api"
	"nhs-bank-notifier/pkg/config"
	"nhs-bank-notifier/pkg/logger"
	"nhs-bank-notifier/pkg/notifier"

	"github.com/jftuga/TtlMap"
)

func main() {
	// Initialize configuration and logging
	cfg := config.LoadConfig()

	logger.Init(cfg.LogLevel)
	log := logger.GetLogger()

	log.Debug("Configuration loaded: ", cfg)

	log.Info("Starting NHS Bank Shift Notifier...")

	// Initialize TtlMap
	ttlMap := TtlMap.New[string](cfg.MaxTTL, 3, time.Second*1, false)
	defer ttlMap.Close()

	// Initialize HTTP client for API calls
	httpClient := api.NewClient()

	// Perform login
	if err := api.Login(httpClient, cfg.LoginURL, cfg.NHSUsername, cfg.NHSPassword); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	for {
		log.Debug("Checking for new shifts...")
		shifts := api.FetchShifts(httpClient, ttlMap, cfg)
		if len(shifts) > 0 {
			message := api.FormatShiftsMessage(shifts)
			log.Debugf("New shifts available: %s", message)
			if err := notifier.SendTelegramMessage(cfg.TelegramToken, cfg.TelegramChatID, message); err != nil {
				log.Errorf("Error sending Telegram message: %v", err)
			}
		} else {
			log.Debug("No new shifts available.")
		}
		time.Sleep(time.Duration(cfg.CheckIntervalMins) * time.Minute)
	}
}
