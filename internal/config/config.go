package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type LarkConfig struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type WecomConfig struct {
	CorpID string `json:"corp_id"`
	Secret string `json:"secret"`
}

type DingtalkConfig struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type WxConfig struct {
	Token string `json:"token"`
}

type Config struct {
	Lark     LarkConfig     `json:"lark"`
	Wecom    WecomConfig    `json:"wecom"`
	Dingtalk DingtalkConfig `json:"dingtalk"`
	Wx       WxConfig       `json:"wx"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	path := filepath.Join(os.Getenv("HOME"), ".config", "chathub", "config.json")
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	if v := os.Getenv("CHATHUB_LARK_APP_ID"); v != "" {
		cfg.Lark.AppID = v
	}
	if v := os.Getenv("CHATHUB_LARK_APP_SECRET"); v != "" {
		cfg.Lark.AppSecret = v
	}
	if v := os.Getenv("CHATHUB_WECOM_CORP_ID"); v != "" {
		cfg.Wecom.CorpID = v
	}
	if v := os.Getenv("CHATHUB_WECOM_SECRET"); v != "" {
		cfg.Wecom.Secret = v
	}
	if v := os.Getenv("CHATHUB_DINGTALK_APP_KEY"); v != "" {
		cfg.Dingtalk.AppKey = v
	}
	if v := os.Getenv("CHATHUB_DINGTALK_APP_SECRET"); v != "" {
		cfg.Dingtalk.AppSecret = v
	}
	if v := os.Getenv("CHATHUB_WX_TOKEN"); v != "" {
		cfg.Wx.Token = v
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path := filepath.Join(os.Getenv("HOME"), ".config", "chathub", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
