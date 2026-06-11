package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type RiskConfig map[string]string

type LargeFilesConfig struct {
	MinSizeMB    int      `mapstructure:"min_size_mb"`
	ExcludePaths []string `mapstructure:"exclude_paths"`
}

type DuplicatesConfig struct {
	MinSizeMB         int      `mapstructure:"min_size_mb"`
	ExcludeExtensions []string `mapstructure:"exclude_extensions"`
}

type ReportConfig struct {
	SaveDir   string `mapstructure:"save_dir"`
	KeepCount int    `mapstructure:"keep_count"`
}

type Config struct {
	Risk       RiskConfig       `mapstructure:"risk"`
	LargeFiles LargeFilesConfig `mapstructure:"large_files"`
	Duplicates DuplicatesConfig `mapstructure:"duplicates"`
	Report     ReportConfig     `mapstructure:"report"`
}

func DefaultConfig() *Config {
	return &Config{
		Risk: RiskConfig{},
		LargeFiles: LargeFilesConfig{
			MinSizeMB:    100,
			ExcludePaths: []string{},
		},
		Duplicates: DuplicatesConfig{
			MinSizeMB:         1,
			ExcludeExtensions: []string{".exe", ".dll", ".sys"},
		},
		Report: ReportConfig{
			SaveDir:   "~/.slimming/reports",
			KeepCount: 30,
		},
	}
}

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if configPath := os.Getenv("SLIMMING_CONFIG"); configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(homeDir, ".slimming"))
		}
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
