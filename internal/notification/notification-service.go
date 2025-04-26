package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/khorzhenwin/go-chujang/internal/models"
	"net/http"
)

type Service struct {
	notificationConfig models.TelegramNotifier
}

func NewService(notificationConfig *models.TelegramNotifier) *Service {
	return &Service{notificationConfig: *notificationConfig}
}

func Send(message string, t *models.TelegramNotifier) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := map[string]string{
		"chat_id": t.ChatID,
		"text":    message,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API responded with status: %d", resp.StatusCode)
	}
	return nil
}
