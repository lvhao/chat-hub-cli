package cmd

import (
	"fmt"

	"github.com/lvhao/chat-hub-cli/internal/config"
	"github.com/lvhao/chat-hub-cli/internal/output"
	"github.com/lvhao/chat-hub-cli/internal/platform/wx"
	"github.com/spf13/cobra"
)

var wxCmd = &cobra.Command{Use: "wx", Short: "WeChat (微信) commands"}

var wxSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message via WeChat",
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
		c := wx.New(cfg.Wx.Token)
		if err := c.SendMessage(to, text); err != nil {
			output.Failure(err)
		}
		output.Success(map[string]string{"status": "sent"})
		return nil
	},
}

var wxReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages from WeChat",
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
		c := wx.New(cfg.Wx.Token)
		msgs, err := c.ReadMessages(from, limit)
		if err != nil {
			output.Failure(err)
		}
		output.Success(msgs)
		return nil
	},
}

func init() {
	wxSendCmd.Flags().String("to", "", "Recipient ID")
	wxSendCmd.Flags().String("text", "", "Message text")
	wxReadCmd.Flags().String("from", "", "Chat/user ID")
	wxReadCmd.Flags().Int("limit", 10, "Max messages to fetch")
	wxCmd.AddCommand(wxSendCmd, wxReadCmd)
	rootCmd.AddCommand(wxCmd)
}
