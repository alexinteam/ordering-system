package cmd

import (
	"log"

	"billing-service/internal/server"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the billing service server",
	Long: `Start the billing service HTTP server with the specified configuration.
The server will handle billing operations including account management and payment processing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start server
		srv := server.NewServer()
		if err := srv.Run(":8081"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
