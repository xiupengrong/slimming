package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看/编辑配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败: %w", err)
		}

		configDir := filepath.Join(homeDir, ".slimming")
		configFile := filepath.Join(configDir, "config.toml")

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("配置文件不存在: %s\n", configFile)
			fmt.Println("使用默认配置运行")
			return nil
		}

		fmt.Printf("配置文件: %s\n", configFile)

		content, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}

		fmt.Printf("\n%s\n", string(content))

		return nil
	},
}
