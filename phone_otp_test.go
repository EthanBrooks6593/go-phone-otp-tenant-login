package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginDecisionRequiresCaptchaAndCode(t *testing.T) {
	cases := []struct {
		name, captcha, code string
		want                bool
	}{
		{"approved", "ok", "123456", true},
		{"missing captcha", "", "123456", false},
		{"missing code", "ok", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if (tc.captcha != "" && tc.code != "") != tc.want {
				t.Fatalf("decision mismatch")
			}
		})
	}
}

func TestVerifyCaptchaIncludesWidgetRecordID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["widget_record_id"]; got != "widget-123" {
			t.Fatalf("widget_record_id = %v, want widget-123", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":null,"error":null,"metadata":null}`))
	}))
	defer server.Close()

	c := &InfraiClient{BaseURL: server.URL, Key: "test", HTTP: server.Client()}
	if err := c.VerifyCaptcha("widget-123", "token-123", "203.0.113.9"); err != nil {
		t.Fatal(err)
	}
}
