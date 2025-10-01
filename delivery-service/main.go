package main

import (
	"log"
	"os"

	"delivery-service/internal/server"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "delivery-service",
		Short: "Delivery service for order system",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the delivery service server",
		Run: func(cmd *cobra.Command, args []string) {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8084"
			}

			srv := server.NewServer()
			log.Printf("Starting delivery service on port %s", port)
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
