package cmd

import (
	"antoine-cli/internal/config"
	"fmt"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Antoine configuration",
	Long:  `Configure API keys, preferences, and settings for optimal Antoine experience.`,
}

var trendsCmd = &cobra.Command{
	Use:   "trends",
	Short: "Show current tech trends",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current tech trends: AI, Blockchain, Web3, IoT, Edge Computing")
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.Show()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		config.Set(key, value)
		fmt.Printf("✅ Set %s = %s\n", key, value)
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
}
