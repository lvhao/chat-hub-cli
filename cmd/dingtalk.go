package cmd

import (
	"fmt"

	"github.com/lvhao/chat-hub-cli/internal/config"
	"github.com/lvhao/chat-hub-cli/internal/output"
	"github.com/lvhao/chat-hub-cli/internal/platform/dingtalk"
	"github.com/spf13/cobra"
)

var dingtalkCmd = &cobra.Command{Use: "dingtalk", Short: "DingTalk (钉钉) commands"}

var dingtalkSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message via DingTalk",
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
		c := dingtalk.New(cfg.Dingtalk.AppKey, cfg.Dingtalk.AppSecret)
		if err := c.SendMessage(to, text); err != nil {
			output.Failure(err)
		}
		output.Success(map[string]string{"status": "sent"})
		return nil
	},
}

var dingtalkReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages from DingTalk",
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
		c := dingtalk.New(cfg.Dingtalk.AppKey, cfg.Dingtalk.AppSecret)
		msgs, err := c.ReadMessages(from, limit)
		if err != nil {
			output.Failure(err)
		}
		output.Success(msgs)
		return nil
	},
}

func init() {
	dingtalkSendCmd.Flags().String("to", "", "Recipient ID")
	dingtalkSendCmd.Flags().String("text", "", "Message text")
	dingtalkReadCmd.Flags().String("from", "", "Chat/user ID")
	dingtalkReadCmd.Flags().Int("limit", 10, "Max messages to fetch")
	dingtalkCmd.AddCommand(dingtalkSendCmd, dingtalkReadCmd)
	rootCmd.AddCommand(dingtalkCmd)
}
