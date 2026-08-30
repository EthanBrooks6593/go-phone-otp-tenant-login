package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type loginRequest struct{ Phone, WidgetRecordID, CaptchaToken, CaptchaIP, Code string }

func login(c *InfraiClient, in loginRequest) (map[string]any, error) {
	if in.WidgetRecordID == "" || in.CaptchaToken == "" || in.Code == "" {
		return nil, fmt.Errorf("widget record ID, captcha, and code are required")
	}
	if err := c.VerifyCaptcha(in.WidgetRecordID, in.CaptchaToken, in.CaptchaIP); err != nil {
		return nil, err
	}
	if err := c.SendCode(in.Phone); err != nil {
		return nil, err
	}
	return c.VerifyCode(in.Phone, in.Code)
}

func main() {
	c, err := NewClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in loginRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		data, err := login(c, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "data": data})
	})
	fmt.Println("listening on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
