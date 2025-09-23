package cmd

import (
	"fmt"
	"log"

	"billing-service/internal/config"
	"billing-service/internal/database"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		db, err := database.Connect(cfg.Database)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		if err := database.Migrate(db); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}

		fmt.Println("Database migrations completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
