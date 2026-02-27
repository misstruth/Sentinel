package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Config Webhook 配置
type Config struct {
	ID      uint
	Name    string
	URL     string
	Secret  string
	Enabled bool
}

// Payload 通知载荷
type Payload struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
	Time  string      `json:"time"`
}

// Client Webhook 客户端
type Client struct {
	config *Config
	http   *http.Client
}

// NewClient 创建客户端
func NewClient(cfg *Config) *Client {
	return &Client{
		config: cfg,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send 发送通知
func (c *Client) Send(ctx context.Context, p *Payload) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.URL, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.Secret != "" {
		req.Header.Set("X-Webhook-Secret", c.config.Secret)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
