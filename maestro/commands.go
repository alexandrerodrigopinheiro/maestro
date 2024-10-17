package maestro

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Command represents a command to be executed.
type Command interface {
	Execute(args []string)
}

// Dependency Management Commands (Composer-like)

// NewProjectCommand creates a new Jazz project with a default structure.
type NewProjectCommand struct{}

// Execute initializes a new project with the provided name.
func (c *NewProjectCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a project name.")
		return
	}
	InitializeProject(args[0])
}

// InstallCommand installs backend and frontend dependencies.
type InstallCommand struct{}

// Execute installs dependencies for Go and React.
func (c *InstallCommand) Execute(args []string) {
	fmt.Println("Installing backend dependencies...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		fmt.Printf("Failed to install backend dependencies: %s\n", err)
		return
	}

	fmt.Println("Installing frontend dependencies...")
	if err := runCommand("npm", "install", "--prefix", "frontend"); err != nil {
		fmt.Printf("Failed to install frontend dependencies: %s\n", err)
		return
	}

	fmt.Println("Dependencies installed successfully!")
}

// AddDependencyCommand adds a new dependency to backend or frontend.
type AddDependencyCommand struct{}

// Execute adds a new dependency to the specified environment.
func (c *AddDependencyCommand) Execute(args []string) {
	if len(args) < 2 {
		fmt.Println("Please specify the environment (go/npm) and the package name.")
		return
	}

	environment := args[0]
	packageName := args[1]

	switch environment {
	case "go":
		fmt.Printf("Adding Go package: %s\n", packageName)
		if err := runCommand("go", "get", packageName); err != nil {
			fmt.Printf("Failed to add Go package: %s\n", err)
			return
		}
	case "npm":
		fmt.Printf("Adding npm package: %s\n", packageName)
		if err := runCommand("npm", "install", packageName, "--prefix", "frontend"); err != nil {
			fmt.Printf("Failed to add npm package: %s\n", err)
			return
		}
	default:
		fmt.Println("Unknown environment. Use 'go' or 'npm'.")
		return
	}

	fmt.Println("Dependency added successfully!")
}

// MigrateCommand runs the database migrations.
type MigrateCommand struct{}

// Execute runs the migration script.
func (c *MigrateCommand) Execute(args []string) {
	fmt.Println("Running database migrations...")
	if err := runCommand("go", "run", "migrations/main.go"); err != nil {
		fmt.Printf("Migration failed: %s\n", err)
		return
	}
	fmt.Println("Migrations completed successfully!")
}

// ServeCommand starts the development server.
type ServeCommand struct{}

// Execute starts the development server for both backend and frontend.
func (c *ServeCommand) Execute(args []string) {
	host := "localhost"
	port := "8080"

	// Check if host and port are provided as arguments
	if len(args) >= 1 {
		host = args[0]
	}
	if len(args) >= 2 {
		port = args[1]
	}

	fmt.Printf("Starting the development server on %s:%s...\n", host, port)

	// Start Backend Server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Backend server is running!"))
		})
		backendAddress := fmt.Sprintf("%s:%s", host, port)
		fmt.Printf("Backend server is running at http://%s\n", backendAddress)
		if err := http.ListenAndServe(backendAddress, mux); err != nil {
			log.Fatalf("Failed to start backend server: %s\n", err)
		}
	}()

	// Start Frontend Dev Server
	if err := runCommand("npm", "start", "--prefix", "frontend"); err != nil {
		fmt.Printf("Failed to start frontend development server: %s\n", err)
		return
	}

	fmt.Println("Server started successfully!")
}

// MakeModelCommand creates a new model file.
type MakeModelCommand struct{}

// Execute creates a new model file with GORM support.
func (c *MakeModelCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a model name.")
		return
	}

	modelName := args[0]
	filename := fmt.Sprintf("backend/models/%s.go", modelName)

	content := fmt.Sprintf(`package models

import (
    "gorm.io/gorm"
)

type %s struct {
    gorm.Model
    Name  string
    Email string
}`, modelName)

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to create model file: %s\n", err)
		return
	}

	fmt.Printf("Model file created: %s\n", filename)
}

// NewMigrationCommand creates a new migration file.
type NewMigrationCommand struct{}

// Execute creates a new migration file with the provided name.
func (c *NewMigrationCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a migration name.")
		return
	}

	migrationName := args[0]
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("backend/migrations/%s_%s.go", timestamp, migrationName)

	content := fmt.Sprintf(`package migrations

import (
	"database/sql"
	"fmt"
)

// Up is executed when this migration is applied
func Up(db *sql.DB) {
	fmt.Println("Applying migration: %s")
	// Add migration logic here
}

// Down is executed when this migration is reverted
func Down(db *sql.DB) {
	fmt.Println("Reverting migration: %s")
	// Add revert logic here
}
`, migrationName, migrationName)

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to create migration file: %s\n", err)
		return
	}

	fmt.Printf("Migration file created: %s\n", filename)
}

// MakeSchemaCommand creates tables in the database using GORM.
type MakeSchemaCommand struct{}

// Execute initializes the database schema using GORM.
func (c *MakeSchemaCommand) Execute(args []string) {
	fmt.Println("Creating database schema using GORM...")

	// Open a connection to the SQLite database (can be replaced with other DB drivers)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to the database: %s\n", err)
		return
	}

	// Define models
	type User struct {
		ID    uint `gorm:"primaryKey"`
		Name  string
		Email string
	}

	// Automatically migrate the schema
	if err := db.AutoMigrate(&User{}); err != nil {
		fmt.Printf("Failed to migrate schema: %s\n", err)
		return
	}

	fmt.Println("Database schema created successfully!")
}
