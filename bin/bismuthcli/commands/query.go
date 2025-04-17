package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewQueryCommand() *cobra.Command {
	var (
		serverURL string
		query     string
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Execute MDX/XMLA queries",
	}

	runCmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute query",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Executing query on %s:\n%s\n", serverURL, query)
			// Здесь логика выполнения запроса
		},
	}

	runCmd.Flags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "Server URL")
	runCmd.Flags().StringVarP(&query, "query", "q", "", "Query to execute")
	_ = runCmd.MarkFlagRequired("query")

	cmd.AddCommand(runCmd)
	return cmd
}
