package cmd

import (
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "扫描可清理内容",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
