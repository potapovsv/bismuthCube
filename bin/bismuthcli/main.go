package main

import (
	"github.com/potapovsv/bismuthCube/bin/bismuthcli/commands"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bismuthcli",
		Short: "BismuthCube management CLI",
		Long:  `Command line interface for BismuthCube OLAP server`,
	}

	// Добавляем подкоманды
	rootCmd.AddCommand(commands.NewServerCommand())
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewQueryCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
