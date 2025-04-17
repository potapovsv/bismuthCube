package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Server management",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "start",
			Short: "Start the server",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Starting BismuthCube server...")
				// Здесь вызов основной логики сервера
			},
		},
		&cobra.Command{
			Use:   "stop",
			Short: "Stop the server",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Stopping server...")
			},
		},
	)

	return cmd
}
