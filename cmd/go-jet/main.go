package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/sqlite"
	"github.com/go-jet/jet/v2/generator/template"
	sqlitelib "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

func main() {
	app := &cli.Command{
		Name:  "go-jet",
		Usage: "Generate go-jet code with custom types for claw project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "database",
				Aliases:  []string{"db"},
				Usage:    "Database connection string",
				Sources:  cli.EnvVars("GOOSE_DBSTRING"),
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output directory for generated code",
				Value:   "lib/claw/gen",
			},
			&cli.StringFlag{
				Name:    "package",
				Aliases: []string{"p"},
				Usage:   "Package name for generated code",
				Value:   "gen",
			},
		},
		Action: generateCode,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateCode(ctx context.Context, c *cli.Command) error {
	dbPath := c.String("database")
	outputDir := c.String("output")
	packageName := c.String("package")

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

	// Create output directory
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Generating go-jet code from database: %s\n", dbPath)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("Package name: %s\n", packageName)

	// Generate code using go-jet with custom template configuration
	err = sqlite.GenerateDB(db,
		outputDir,
		template.Default(sqlitelib.Dialect).
			UseSchema(func(schemaMetadata metadata.Schema) template.Schema {
				return template.DefaultSchema(schemaMetadata).
					UseModel(template.DefaultModel().
						UseTable(func(tableMetadata metadata.Table) template.TableModel {
							// Skip goose version table
							if tableMetadata.Name == "goose_db_version" {
								tableModel := template.DefaultTableModel(tableMetadata)
								tableModel.Skip = true
								return tableModel
							}

							return template.DefaultTableModel(tableMetadata).
								UseField(func(columnMetadata metadata.Column) template.TableModelField {
									defaultField := template.DefaultTableModelField(columnMetadata)

									// Apply custom type mapping based on field naming rules
									customType := getCustomType(columnMetadata.Name, columnMetadata.DataType.Name, columnMetadata.IsNullable)
									if customType != nil {
										defaultField.Type = *customType
									}

									return defaultField
								})
						}))
			}),
	)
	if err != nil {
		return fmt.Errorf("failed to generate go-jet code: %w", err)
	}

	fmt.Println("Go-jet code generation completed successfully!")
	return nil
}

// getCustomType returns the custom type for a field based on naming conventions
func getCustomType(columnName, dataType string, isNullable bool) *template.Type {
	columnName = strings.ToLower(columnName)
	dataType = strings.ToUpper(dataType)

	// Rule 1: Timestamp fields with _at suffix (INTEGER -> UnixMilli)
	if strings.HasSuffix(columnName, "_at") && dataType == "INTEGER" {
		if isNullable {
			customType := template.NewType((*types.UnixMilli)(nil))
			return &customType
		}
		customType := template.NewType(types.UnixMilli{})
		return &customType
	}

	// Rule 2: Duration fields with _dur suffix (INTEGER -> DurationMilli)
	if strings.HasSuffix(columnName, "_dur") && dataType == "INTEGER" {
		if isNullable {
			customType := template.NewType((*types.DurationMilli)(nil))
			return &customType
		}
		customType := template.NewType(types.DurationMilli{})
		return &customType
	}

	// Rule 3: Boolean fields with is_ prefix (INTEGER -> Bool)
	if strings.HasPrefix(columnName, "is_") && dataType == "INTEGER" {
		if isNullable {
			customType := template.NewType((*types.Bool)(nil))
			return &customType
		}
		customType := template.NewType(types.Bool(false))
		return &customType
	}

	// Rule 4: Default INTEGER fields -> int64
	if dataType == "INTEGER" {
		if isNullable {
			customType := template.NewType((*int64)(nil))
			return &customType
		}
		customType := template.NewType(int64(0))
		return &customType
	}

	// Rule 5: Default REAL fields -> float64
	if dataType == "REAL" {
		if isNullable {
			customType := template.NewType((*float64)(nil))
			return &customType
		}
		customType := template.NewType(float64(0))
		return &customType
	}

	// Return nil for default go-jet type mapping
	return nil
}
