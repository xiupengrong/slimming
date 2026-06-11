package cmd

import (
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理C盘空间",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
