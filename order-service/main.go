package main

import (
	"log"
	"order-service/internal/server"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "order-service",
		Short: "Order service for order system",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the order service server",
		Run: func(cmd *cobra.Command, args []string) {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8081"
			}

			srv := server.NewServer()
			log.Printf("Starting order service on port %s", port)
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
