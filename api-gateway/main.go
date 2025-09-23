package main

import (
	"api-gateway/internal/server"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "api-gateway",
		Short: "API Gateway for order system",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the API Gateway server",
		Run: func(cmd *cobra.Command, args []string) {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			srv := server.NewServer()
			log.Printf("Starting API Gateway on port %s", port)
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
