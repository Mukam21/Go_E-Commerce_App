package notification

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Mukam21/Go_E-Commerce_App/config"
)

type NotificationClient interface {
	SendSMS(phone string, message string) error
}

type notificationClient struct {
	config config.AppConfig
}

func (c notificationClient) SendSMS(phone string, message string) error {

	// DEV режим — не шлём SMS реально
	if c.config.Env == "dev" {
		log.Println("SMS (dev):", phone, message)
		return nil
	}

	apiURL := "https://sms.ru/sms/send"

	data := url.Values{}
	data.Set("api_id", c.config.SMSRuApiKey)
	data.Set("to", phone)
	data.Set("msg", message)
	data.Set("json", "1")

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result["status"] != "OK" {
		return errors.New("sms.ru send failed")
	}

	return nil
}

func NewNotificationClient(config config.AppConfig) NotificationClient {
	return &notificationClient{
		config: config,
	}
}
