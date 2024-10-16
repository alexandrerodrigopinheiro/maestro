package maestro

import (
	"fmt"
	"log"
	"net/http"
)

// Command represents a command to be executed.
type Command interface {
	Execute(args []string)
}

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
