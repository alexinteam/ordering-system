package cmd

import (
	"log"

	"api-gateway/internal/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API Gateway server",
	Long:  `Start the API Gateway server.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.NewServer()

		if err := srv.Run(":8080"); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
