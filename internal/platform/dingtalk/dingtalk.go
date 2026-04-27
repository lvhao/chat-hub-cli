package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lvhao/chat-hub-cli/internal/platform"
)

type Client struct {
	appKey    string
	appSecret string
	httpClient *http.Client
}

func New(appKey, appSecret string) *Client {
	return &Client{appKey: appKey, appSecret: appSecret, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) getToken() (string, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", c.appKey, c.appSecret)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		ErrCode     int    `json:"errcode"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("dingtalk auth failed: code %d", result.ErrCode)
	}
	return result.AccessToken, nil
}

func (c *Client) SendMessage(to, text string) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{
		"receiver_id": to,
		"msg": map[string]any{
			"msgtype": "text",
			"text":    map[string]string{"content": text},
		},
	})
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s", token)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct{ ErrCode int `json:"errcode"` }
	json.NewDecoder(resp.Body).Decode(&result)
	if result.ErrCode != 0 {
		return fmt.Errorf("dingtalk send failed: code %d", result.ErrCode)
	}
	return nil
}

func (c *Client) ReadMessages(from string, limit int) ([]platform.Message, error) {
	return []platform.Message{}, nil
}
