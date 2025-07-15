package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/tigorlazuardi/claw/migrations"
	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

func main() {
	app := &cli.Command{
		Name:  "goose",
		Usage: "Database migration tool for development",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database",
				Aliases:  []string{"db"},
				Usage:    "Database connection string",
				Sources:  cli.EnvVars("GOOSE_DBSTRING"),
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "reset",
				Aliases: []string{"r"},
				Usage:   "Reset all migrations and reapply them (WARNING: DATA LOSS)",
				Value:   true, // Default to true for development
			},
			&cli.StringFlag{
				Name:    "fs-path",
				Aliases: []string{"f"},
				Usage:   "Custom filesystem path for migrations (optional)",
			},
		},
		Action: runMigrations,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(ctx context.Context, c *cli.Command) error {
	dbPath := c.String("database")
	reset := c.Bool("reset")
	fsPath := c.String("fs-path")

	// Check if database file exists
	_, err := os.Stat(dbPath)
	dbExists := err == nil

	// If database doesn't exist and reset is true, disable reset
	if !dbExists && reset {
		fmt.Println("Database file doesn't exist, skipping reset operation")
		reset = false
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Prepare migration options
	var options []migrations.MigrateOption

	// Add reset option
	options = append(options, migrations.WithReset(reset))

	// Add custom filesystem if provided
	if fsPath != "" {
		customFS := os.DirFS(fsPath)
		options = append(options, migrations.WithFS(customFS))
	}

	// Run migrations
	fmt.Printf("Running migrations with reset=%v\n", reset)
	if fsPath != "" {
		fmt.Printf("Using custom filesystem path: %s\n", fsPath)
	}

	if err := migrations.Migrate(ctx, db, options...); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations completed successfully!")
	return nil
}
