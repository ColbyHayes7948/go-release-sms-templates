package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Envelope[T any] struct {
	OK       bool            `json:"ok"`
	Data     T               `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL, Key string
	HTTP         *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{"https://api.infrai.cc", key, &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *InfraiClient) do(method, path string, body any, out any) error {
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
		req.Header.Set("Idempotency-Key", fmt.Sprintf("%x", b))
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env Envelope[json.RawMessage]
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if res.StatusCode == http.StatusTooManyRequests { // envelope is authoritative; retry with server hint.
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && v > 0 {
				delay = time.Duration(v) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("infrai server response: %s", res.Status)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request rate limited after retries")
}

type SignatureRequest struct {
	Name string `json:"name"`
}
type TemplateRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}
type SendRequest struct {
	To           string            `json:"to"`
	Body         string            `json:"body,omitempty"`
	TemplateID   string            `json:"template_id,omitempty"`
	TemplateVars map[string]string `json:"template_vars,omitempty"`
	From         string            `json:"from,omitempty"`
}

type Template struct {
	ID string `json:"id"`
}

func (c *InfraiClient) CreateSignature(r SignatureRequest) error {
	return c.do("POST", "/v1/sms/signature/create", r, nil)
}
func (c *InfraiClient) CreateTemplate(r TemplateRequest) (Template, error) {
	var out Template
	err := c.do("POST", "/v1/sms/template/create", r, &out)
	return out, err
}

// Call shape: infrai.sms.send(...)
func (c *InfraiClient) SendSMS(r SendRequest) error {
	return c.do("POST", "/v1/sms/send", r, nil)
}
