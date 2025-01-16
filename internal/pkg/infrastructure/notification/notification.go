package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationService interface {
	SendMessage(message string) error
}

type TelegramNotifier struct {
	APIURL string
}

func NewTelegramNotifier(apiURL string) NotificationService {
	return &TelegramNotifier{APIURL: apiURL}
}

func (t *TelegramNotifier) SendMessage(message string) error {
	payload := map[string]string{
		"message": message,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(t.APIURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}
	return nil
}
