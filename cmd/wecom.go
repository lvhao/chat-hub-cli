package cmd

import (
	"fmt"

	"github.com/lvhao/chat-hub-cli/internal/config"
	"github.com/lvhao/chat-hub-cli/internal/output"
	"github.com/lvhao/chat-hub-cli/internal/platform/wecom"
	"github.com/spf13/cobra"
)

var wecomCmd = &cobra.Command{Use: "wecom", Short: "WeCom (企业微信) commands"}

var wecomSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message via WeCom",
	RunE: func(cmd *cobra.Command, args []string) error {
		to, _ := cmd.Flags().GetString("to")
		text, _ := cmd.Flags().GetString("text")
		if to == "" || text == "" {
			return fmt.Errorf("--to and --text are required")
		}
		cfg, err := config.Load()
		if err != nil {
			output.Failure(err)
		}
		c := wecom.New(cfg.Wecom.CorpID, cfg.Wecom.Secret)
		if err := c.SendMessage(to, text); err != nil {
			output.Failure(err)
		}
		output.Success(map[string]string{"status": "sent"})
		return nil
	},
}

var wecomReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages from WeCom",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		limit, _ := cmd.Flags().GetInt("limit")
		if from == "" {
			return fmt.Errorf("--from is required")
		}
		cfg, err := config.Load()
		if err != nil {
			output.Failure(err)
		}
		c := wecom.New(cfg.Wecom.CorpID, cfg.Wecom.Secret)
		msgs, err := c.ReadMessages(from, limit)
		if err != nil {
			output.Failure(err)
		}
		output.Success(msgs)
		return nil
	},
}

func init() {
	wecomSendCmd.Flags().String("to", "", "Recipient ID")
	wecomSendCmd.Flags().String("text", "", "Message text")
	wecomReadCmd.Flags().String("from", "", "Chat/user ID")
	wecomReadCmd.Flags().Int("limit", 10, "Max messages to fetch")
	wecomCmd.AddCommand(wecomSendCmd, wecomReadCmd)
	rootCmd.AddCommand(wecomCmd)
}
