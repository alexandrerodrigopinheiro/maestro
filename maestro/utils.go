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
func InitializeProject(projectName string) error {
	fmt.Printf("Initializing a new Go project: %s\n", projectName)

	// Create project root directory
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %s", err)
	}

	// Verify if the directory was created successfully
	if _, err := os.Stat(projectName); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist after creation attempt")
	}

	// Navigate to project directory
	if err := os.Chdir(projectName); err != nil {
		return fmt.Errorf("failed to navigate to project directory: %s", err)
	}

	// Initialize Go module
	if err := runCommand("go", "mod", "init", projectName); err != nil {
		return fmt.Errorf("failed to initialize Go module: %s", err)
	}

	// Create initial backend folders similar to a Laravel structure
	folders := []string{
		"backend/cmd",
		"backend/pkg",
		"backend/internal",
		"backend/configs",
		"backend/migrations",
		"backend/routes",
		"backend/controllers",
		"backend/models",
		"backend/middleware",
	}
	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("failed to create folder %s: %s", folder, err)
		}
	}

	return nil
}

// CleanupProject removes the project directory in case of errors during setup.
func CleanupProject(projectName string) {
	fmt.Printf("Cleaning up incomplete project '%s'...\n", projectName)
	if err := os.RemoveAll(projectName); err != nil {
		fmt.Printf("Failed to clean up project directory: %s\n", err)
	} else {
		fmt.Println("Project cleanup completed.")
	}
}
