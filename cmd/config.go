package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看/编辑配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
