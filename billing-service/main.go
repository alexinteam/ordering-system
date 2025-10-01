package main

import (
	"log"
	"os"

	"billing-service/internal/server"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "billing-service",
		Short: "Billing service for order system",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the billing service server",
		Run: func(cmd *cobra.Command, args []string) {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8082"
			}

			srv := server.NewServer()
			log.Printf("Starting billing service on port %s", port)
			if err := srv.Run(":" + port); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
