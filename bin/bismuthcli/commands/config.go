package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("Error reading config: %v\n", err)
				return
			}
			fmt.Println("Current configuration:")
			for k, v := range viper.AllSettings() {
				fmt.Printf("%s: %v\n", k, v)
			}
		},
	}

	showCmd.Flags().StringVarP(&configFile, "config", "c", "config.yml", "Config file path")
	cmd.AddCommand(showCmd)

	return cmd
}
