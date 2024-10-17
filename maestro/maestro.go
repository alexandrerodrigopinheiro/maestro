// maestro/maestro.go

package maestro

import (
	"fmt"
	"os"
)

// commands holds a registry of available commands.
var commands = map[string]Command{
	"new":          &NewProjectCommand{},
	"install":      &InstallCommand{},
	"add":          &AddDependencyCommand{},
	"migrate":      &MigrateCommand{},
	"serve":        &ServeCommand{},
	"make:model":   &MakeModelCommand{},
	"make:migrate": &NewMigrationCommand{},
	"make:schema":  &MakeSchemaCommand{},
}

// Main function to handle commands.
func Main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command (use 'help' to see available commands).")
		return
	}

	commandName := os.Args[1]
	if commandName == "help" {
		listCommands()
		return
	}

	command, exists := commands[commandName]
	if !exists {
		fmt.Println("Unknown command. Use 'help' to see available commands.")
		return
	}

	command.Execute(os.Args[2:])
}

// listCommands prints the available commands with descriptions.
func listCommands() {
	fmt.Println("Available commands:")
	fmt.Println("- new          Creates a new project with a default structure.")
	fmt.Println("- install      Installs backend and frontend dependencies.")
	fmt.Println("- add          Adds a new dependency to backend or frontend.")
	fmt.Println("- migrate      Runs the database migrations.")
	fmt.Println("- serve        Starts the development server for both backend and frontend.")
	fmt.Println("- make:model   Creates a new model file with GORM support.")
	fmt.Println("- make:migrate Creates a new migration file.")
	fmt.Println("- make:schema  Initializes the database schema using GORM.")
}
