package cmd

import (
	"fmt"
	"strings"

	"github.com/lvhao/chat-hub-cli/internal/config"
	"github.com/lvhao/chat-hub-cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{Use: "config", Short: "Manage chathub configuration"}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value (e.g. lark.app_id)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, val := args[0], args[1]
		cfg, err := config.Load()
		if err != nil {
			output.Failure(err)
		}
		switch key {
		case "lark.app_id":
			cfg.Lark.AppID = val
		case "lark.app_secret":
			cfg.Lark.AppSecret = val
		case "wecom.corp_id":
			cfg.Wecom.CorpID = val
		case "wecom.secret":
			cfg.Wecom.Secret = val
		case "dingtalk.app_key":
			cfg.Dingtalk.AppKey = val
		case "dingtalk.app_secret":
			cfg.Dingtalk.AppSecret = val
		case "wx.token":
			cfg.Wx.Token = val
		default:
			return fmt.Errorf("unknown key: %s", key)
		}
		if err := config.Save(cfg); err != nil {
			output.Failure(err)
		}
		output.Success(map[string]string{"key": key, "status": "saved"})
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration (secrets masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			output.Failure(err)
		}
		mask := func(s string) string {
			if s == "" {
				return ""
			}
			return strings.Repeat("*", len(s))
		}
		output.Success(map[string]any{
			"lark":     map[string]string{"app_id": cfg.Lark.AppID, "app_secret": mask(cfg.Lark.AppSecret)},
			"wecom":    map[string]string{"corp_id": cfg.Wecom.CorpID, "secret": mask(cfg.Wecom.Secret)},
			"dingtalk": map[string]string{"app_key": cfg.Dingtalk.AppKey, "app_secret": mask(cfg.Dingtalk.AppSecret)},
			"wx":       map[string]string{"token": mask(cfg.Wx.Token)},
		})
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
}
