// pkg/notifier/telegram.go
package notifier

import (
	"bytes"
	"fmt"
	"net/http"
	"nhs-bank-notifier/pkg/logger"
)

// SendTelegramMessage sends a message using Telegram Bot API
func SendTelegramMessage(apiToken, chatID, message string) error {
	log := logger.GetLogger()

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", apiToken)
	payload := bytes.NewReader([]byte(fmt.Sprintf(`{"chat_id": "%s", "text": "%s"}`, chatID, message)))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %s", resp.Status)
	}

	log.Infof("Message sent to %s", chatID)
	return nil
}
