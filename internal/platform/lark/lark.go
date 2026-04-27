package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lvhao/chathub/internal/platform"
)

type Client struct {
	appID      string
	appSecret  string
	baseURL    string
	httpClient *http.Client
}

func New(appID, appSecret string) *Client {
	return &Client{appID: appID, appSecret: appSecret, baseURL: "https://open.feishu.cn", httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) getToken() (string, error) {
	body, _ := json.Marshal(map[string]string{"app_id": c.appID, "app_secret": c.appSecret})
	resp, err := c.httpClient.Post(c.baseURL+"/open-apis/auth/v3/tenant_access_token/internal", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Code              int    `json:"code"`
		TenantAccessToken string `json:"tenant_access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", fmt.Errorf("lark auth failed: code %d", result.Code)
	}
	return result.TenantAccessToken, nil
}

func (c *Client) SendMessage(to, text string) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{
		"receive_id": to,
		"msg_type":   "text",
		"content":    fmt.Sprintf(`{"text":"%s"}`, text),
	})
	req, _ := http.NewRequest("POST", c.baseURL+"/open-apis/im/v1/messages?receive_id_type=open_id", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct{ Code int `json:"code"` }
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Code != 0 {
		return fmt.Errorf("lark send failed: code %d", result.Code)
	}
	return nil
}

func (c *Client) ReadMessages(from string, limit int) ([]platform.Message, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/open-apis/im/v1/messages?container_id_type=chat&container_id=%s&page_size=%d", c.baseURL, from, limit)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Code int `json:"code"`
		Data struct {
			Items []struct {
				MessageID  string `json:"message_id"`
				Sender     struct{ ID string `json:"id"` } `json:"sender"`
				Body       struct{ Content string `json:"content"` } `json:"body"`
				CreateTime string `json:"create_time"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("lark read failed: code %d", result.Code)
	}
	msgs := make([]platform.Message, 0, len(result.Data.Items))
	for _, item := range result.Data.Items {
		msgs = append(msgs, platform.Message{
			ID:   item.MessageID,
			From: item.Sender.ID,
			Text: item.Body.Content,
		})
	}
	return msgs, nil
}
