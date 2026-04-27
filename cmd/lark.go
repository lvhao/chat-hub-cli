package cmd

import (
	"fmt"

	"github.com/lvhao/chathub/internal/config"
	"github.com/lvhao/chathub/internal/output"
	"github.com/lvhao/chathub/internal/platform/lark"
	"github.com/spf13/cobra"
)

var larkCmd = &cobra.Command{Use: "lark", Short: "Feishu (Lark) commands"}

var larkSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message via Lark",
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
		c := lark.New(cfg.Lark.AppID, cfg.Lark.AppSecret)
		if err := c.SendMessage(to, text); err != nil {
			output.Failure(err)
		}
		output.Success(map[string]string{"status": "sent"})
		return nil
	},
}

var larkReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages from Lark",
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
		c := lark.New(cfg.Lark.AppID, cfg.Lark.AppSecret)
		msgs, err := c.ReadMessages(from, limit)
		if err != nil {
			output.Failure(err)
		}
		output.Success(msgs)
		return nil
	},
}

func init() {
	larkSendCmd.Flags().String("to", "", "Recipient ID")
	larkSendCmd.Flags().String("text", "", "Message text")
	larkReadCmd.Flags().String("from", "", "Chat/user ID")
	larkReadCmd.Flags().Int("limit", 10, "Max messages to fetch")
	larkCmd.AddCommand(larkSendCmd, larkReadCmd)
	rootCmd.AddCommand(larkCmd)
}
