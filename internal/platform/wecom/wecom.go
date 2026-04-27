package wecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lvhao/chathub/internal/platform"
)

type Client struct {
	corpID     string
	secret     string
	httpClient *http.Client
}

func New(corpID, secret string) *Client {
	return &Client{corpID: corpID, secret: secret, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) getToken() (string, error) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.corpID, c.secret)
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
		return "", fmt.Errorf("wecom auth failed: code %d", result.ErrCode)
	}
	return result.AccessToken, nil
}

func (c *Client) SendMessage(to, text string) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]any{
		"touser":  to,
		"msgtype": "text",
		"agentid": 0,
		"text":    map[string]string{"content": text},
	})
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct{ ErrCode int `json:"errcode"` }
	json.NewDecoder(resp.Body).Decode(&result)
	if result.ErrCode != 0 {
		return fmt.Errorf("wecom send failed: code %d", result.ErrCode)
	}
	return nil
}

func (c *Client) ReadMessages(from string, limit int) ([]platform.Message, error) {
	// 企微消息接收需要配置回调服务器，此处返回空列表作为占位
	return []platform.Message{}, nil
}
