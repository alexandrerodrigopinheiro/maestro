package maestro

import (
	"fmt"
	"os"
	"os/exec"
)

// runCommand runs a shell command and prints its output.
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InitializeProject sets up the initial folder structure and initializes a Go module.
func InitializeProject(projectName string) {
	if projectName == "" {
		fmt.Println("Please provide a project name.")
		return
	}

	fmt.Printf("Initializing a new Go project: %s\n", projectName)
	// Create project root directory
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		fmt.Printf("Failed to create project directory: %s\n", err)
		return
	}

	// Navigate to project directory
	if err := os.Chdir(projectName); err != nil {
		fmt.Printf("Failed to navigate to project directory: %s\n", err)
		return
	}

	// Initialize Go module
	if err := runCommand("go", "mod", "init", projectName); err != nil {
		fmt.Printf("Failed to initialize Go module: %s\n", err)
		return
	}

	// Create initial folders similar to a Laravel structure for Go backend and React frontend
	folders := []string{
		// Backend (Go) folders
		"backend/cmd",
		"backend/pkg",
		"backend/internal",
		"backend/configs",
		"backend/migrations",
		"backend/routes",
		"backend/controllers",
		"backend/models",
		"backend/middleware",
		// Frontend (React TypeScript) folders
		"frontend/public",
		"frontend/src/components",
		"frontend/src/pages",
		"frontend/src/services",
		"frontend/src/styles",
		"frontend/src/utils",
		"frontend/src/assets",
	}
	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			fmt.Printf("Failed to create folder %s: %s\n", folder, err)
			return
		}
	}

	fmt.Println("Project initialized successfully!")
}
