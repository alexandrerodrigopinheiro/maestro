package maestro

import (
	"fmt"
	"os"
)

// commands holds a registry of available commands.
var commands = map[string]Command{
	"new":     &NewProjectCommand{},
	"migrate": &MigrateCommand{},
	"serve":   &ServeCommand{},
}

// Main function to handle commands.
func Main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command (new, migrate, serve).")
		return
	}

	commandName := os.Args[1]
	command, exists := commands[commandName]
	if !exists {
		fmt.Println("Unknown command. Available commands: new, migrate, serve.")
		return
	}

	command.Execute(os.Args[2:])
}
