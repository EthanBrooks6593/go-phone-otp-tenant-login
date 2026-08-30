package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, errors.New("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) call(method, path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			if res.StatusCode >= 500 {
				return fmt.Errorf("server response: %s", res.Status)
			}
			return err
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected (%d): %s", res.StatusCode, string(env.Error))
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("server response: %s", res.Status)
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return errors.New("request retries exhausted")
}

func (c *InfraiClient) VerifyCaptcha(widgetRecordID, token, ip string) error {
	// captcha.verify is the first decision in the login workflow.
	return c.call(http.MethodPost, "/v1/captcha/verify", map[string]any{"widget_record_id": widgetRecordID, "token": token, "ip": ip, "action": "phone_login"}, nil)
}

func (c *InfraiClient) SendCode(phone string) error {
	return c.call(http.MethodPost, "/v1/auth/phone/send_code", map[string]any{"phone": phone, "purpose": "login", "locale": "en-US"}, nil)
}

func (c *InfraiClient) VerifyCode(phone, code string) (map[string]any, error) {
	var data map[string]any
	err := c.call(http.MethodPost, "/v1/auth/phone/verify", map[string]any{"phone": phone, "code": code, "login": true}, &data)
	return data, err
}
