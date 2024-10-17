package maestro

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"time"

	"github.com/alexandrerodrigopinheiro/maestro/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Command represents a command to be executed.
type Command interface {
	Execute(args []string)
}

// Dependency Management Commands (Composer-like)

// Execute initializes a new project with the provided name.
// NewProjectCommand creates a new Jazz project with a default structure.
type NewProjectCommand struct{}

// Execute initializes a new project with the provided name and installs dependencies.
func (c *NewProjectCommand) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a project name.")
		return
	}
	projectName := args[0]

	// Check if project directory already exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		fmt.Printf("A project named '%s' already exists. Aborting.\n", projectName)
		return
	}

	// Initialize project folder structure
	if err := InitializeProject(projectName); err != nil {
		fmt.Printf("Failed to initialize project: %s\n", err)
		CleanupProject(projectName)
		return
	}

	// Install Go dependencies
	fmt.Println("Installing backend dependencies...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		fmt.Printf("Failed to install backend dependencies: %s\n", err)
		CleanupProject(projectName)
		return
	}

	// Initialize React frontend
	fmt.Println("Initializing frontend...")
	if err := runCommand("npx", "create-react-app", "frontend"); err != nil {
		fmt.Printf("Failed to initialize React frontend: %s\n", err)
		CleanupProject(projectName)
		return
	}

	// Install React dependencies
	fmt.Println("Installing frontend dependencies...")
	if err := runCommand("npm", "install", "--prefix", "frontend"); err != nil {
		fmt.Printf("Failed to install frontend dependencies: %s\n", err)
		CleanupProject(projectName)
		return
	}

	// Create additional frontend folders if they don't exist
	additionalFrontendFolders := []string{
		"frontend/src/services",
		"frontend/src/styles",
		"frontend/src/utils",
		"frontend/src/assets",
	}
	for _, folder := range additionalFrontendFolders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			if err := os.MkdirAll(folder, 0755); err != nil {
				fmt.Printf("Failed to create folder %s: %s\n", folder, err)
				CleanupProject(projectName)
				return
			}
		}
	}

	// Create .env file for backend
	fmt.Println("Creating backend .env file...")
	envContentBackend := `# Environment Configuration
APP_NAME=` + projectName + `
APP_ENV=development
APP_HOST=localhost
APP_PORT=8000
API_PORT=8001
APP_KEY=base64:3bH0U5PBJEuJi3vIoBjQzFNEykjKeftLoMRt1+juh38=
APP_DEBUG=true
APP_URL=http://localhost
APP_VERSION=1.3.0
APP_TIMEZONE="America/Sao_Paulo"

LOG_CHANNEL=stack
LOG_DEPRECATIONS_CHANNEL=null
LOG_LEVEL=debug

DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=

REDIS_HOST=127.0.0.1
REDIS_PASSWORD=null
REDIS_PORT=6379
`
	envPathBackend := ".env"
	err := os.WriteFile(envPathBackend, []byte(envContentBackend), 0644)
	if err != nil {
		fmt.Printf("Failed to create backend .env file: %s\n", err)
		CleanupProject(projectName)
		return
	} else {
		fmt.Println("Backend .env file created successfully.")
	}

	// Create .env file for frontend (React)
	fmt.Println("Creating frontend .env file...")
	envContentFrontend := `# React Environment Configuration
REACT_APP_NAME=` + projectName + `
REACT_APP_ENV=development
REACT_APP_API_URL=http://localhost:8000
REACT_APP_VERSION=1.3.0
REACT_APP_DEBUG=true
`
	envPathFrontend := "frontend/.env"
	err = os.WriteFile(envPathFrontend, []byte(envContentFrontend), 0644)
	if err != nil {
		fmt.Printf("Failed to create frontend .env file: %s\n", err)
		CleanupProject(projectName)
		return
	} else {
		fmt.Println("Frontend .env file created successfully.")
	}

	fmt.Println("Project setup completed successfully!")
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

// Execute runs all migration scripts in the migrations folder.
func (c *MigrateCommand) Execute(args []string) {
	// Load environment variables
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	// Set up database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s\n", err)
	}

	// Iterate over migration files in the migrations folder
	migrationPath := "backend/migrations"
	err = filepath.Walk(migrationPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			fmt.Printf("Applying migration: %s\n", info.Name())

			// Load migration plugin
			p, err := plugin.Open(path)
			if err != nil {
				return fmt.Errorf("failed to load migration plugin: %s", err)
			}

			// Look for the Up function in the plugin
			upFunc, err := p.Lookup("Up")
			if err != nil {
				return fmt.Errorf("failed to find 'Up' function in migration %s: %s", path, err)
			}

			// Assert that Up is a function of the right signature and call it
			if up, ok := upFunc.(func(*gorm.DB)); ok {
				up(db)
			} else {
				return fmt.Errorf("invalid 'Up' function signature in migration %s", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to apply migrations: %s\n", err)
	}

	fmt.Println("Migrations completed successfully!")
}

// ServeCommand starts the development server.
type ServeCommand struct{}

// Execute starts the development server for both backend and frontend.
func (c *ServeCommand) Execute(args []string) {
	host := "localhost"
	appPort := "8080" // Porta padrão do React (Frontend)
	apiPort := "8001" // Porta do backend (API), calculada automaticamente

	// Load environment from .env file located in the current directory
	if err := utils.LoadEnv(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	// Get values from environment variables
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production" // Default environment
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		appPort = port
	}
	if port := os.Getenv("API_PORT"); port != "" {
		apiPort = port
	}

	if env == "development" {
		fmt.Printf("Starting the development server on %s:%s (frontend) and %s:%s (backend)...\n", host, appPort, host, apiPort)
	}

	// Start Backend Server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Backend server is running!"))
		})
		backendAddress := fmt.Sprintf("%s:%s", host, apiPort)
		fmt.Printf("Backend server is running at http://%s\n", backendAddress)
		if err := http.ListenAndServe(backendAddress, mux); err != nil {
			log.Fatalf("Failed to start backend server: %s\n", err)
		}
	}()

	// Check if the frontend directory exists
	if _, err := os.Stat("frontend"); os.IsNotExist(err) {
		fmt.Println("Frontend directory not found. Please ensure the frontend project is set up correctly.")
		return
	}

	// Check if package.json exists
	if _, err := os.Stat("frontend/package.json"); os.IsNotExist(err) {
		fmt.Println("package.json not found in the frontend directory. Please initialize a React project in the 'frontend' folder.")
		return
	}

	// Start Frontend Dev Server with specified port
	fmt.Println("Starting frontend server...")
	if err := runCommand("npm", "start", "--prefix", "frontend", "--", "--port", appPort); err != nil {
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
	snakeCaseModelName := utils.ToSnakeCase(modelName)
	pascalCaseModelName := utils.ToPascalCase(snakeCaseModelName)
	filename := fmt.Sprintf("backend/models/%s.go", snakeCaseModelName)

	content := fmt.Sprintf(`package models

import (
    "gorm.io/gorm"
)

type %s struct {
    gorm.Model
}`, pascalCaseModelName)

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

	modelName := args[0]
	pascalCaseModelName := utils.ToPascalCase(modelName)
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("backend/migrations/%s_create_%s_table.go", timestamp, utils.ToSnakeCase(modelName))

	content := fmt.Sprintf(`package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"backend/models"
)

// Up is executed when this migration is applied
func Up(db *gorm.DB) {
	fmt.Println("Applying migration: create %s table")
	// Automigrate the model
	db.AutoMigrate(&models.%s{})
}

// Down is executed when this migration is reverted
func Down(db *gorm.DB) {
	fmt.Println("Reverting migration: drop %s table")
	// Drop the table associated with the model
	db.Migrator().DropTable(&models.%s{})
}
`, utils.ToSnakeCase(modelName), pascalCaseModelName, utils.ToSnakeCase(modelName), pascalCaseModelName)

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to create migration file: %s\n", err)
		return
	}

	fmt.Printf("Migration file created: %s\n", filename)
}
