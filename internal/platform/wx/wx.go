package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lvhao/chat-hub-cli/internal/platform"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func New(token string) *Client {
	return &Client{token: token, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) SendMessage(to, text string) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", c.token)
	body, _ := json.Marshal(map[string]any{
		"touser":  to,
		"msgtype": "text",
		"text":    map[string]string{"content": text},
	})
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct{ ErrCode int `json:"errcode"` }
	json.NewDecoder(resp.Body).Decode(&result)
	if result.ErrCode != 0 {
		return fmt.Errorf("wx send failed: code %d", result.ErrCode)
	}
	return nil
}

func (c *Client) ReadMessages(from string, limit int) ([]platform.Message, error) {
	return []platform.Message{}, nil
}
