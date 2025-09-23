package main

import (
	"billing-service/internal/server"
	"log"
	"os"

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
				port = "8081"
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
