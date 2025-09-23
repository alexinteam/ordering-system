package main

import (
	"log"
	"notification-service/internal/server"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "notification-service",
		Short: "Notification service for order system",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the notification service server",
		Run: func(cmd *cobra.Command, args []string) {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8082"
			}

			srv := server.NewServer()
			log.Printf("Starting notification service on port %s", port)
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
