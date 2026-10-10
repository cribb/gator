package main

import (
	// "bufio"
	// "encoding/json"
	// "strings"
	"fmt"
	"os"

	"github.com/cribb/gator/internal/config"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*cliConfig, []string, error)
}

// map[string]func(*state, command) error

type state struct {
	cfg *config.Config
}

func processCommand(commandLine []string) {
	commands := getCommands()

	// process input (REPL-like)

	commandName := commandLine[0]
	commandArgs := commandLine[1:]

	command, ok := commands[commandName]
	if !ok {
		fmt.Println("Unknown command.")
		cmdHelp(state, commandArgs)
		cmdExit(state, commandArgs)
	}

	err := command.callback(state, commandArgs)
}
func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit gator",
			callback:    cmdExit,
		},
		"quit": {
			name:        "quit",
			description: "Exit gator",
			callback:    cmdExit,
		},
		"help": {
			name:        "help",
			description: "Displays this command list",
			callback:    cmdHelp,
		},
		"login": {
			name:        "login",
			description: "Displays this command list",
			callback:    cmdHelp,
		},
		"register": {
			name:        "register",
			description: "Displays this command list",
			callback:    cmdHelp,
		},
		"users": {
			name:        "users",
			description: "Displays this command list",
			callback:    cmdHelp,
		},
	}
	return commands
}

func cmdExit(config *cliConfig, cmdArgs []string) error {
	fmt.Printf("Gone from the Gator, bye!\n\n")
	os.Exit(0)
	return nil
}

func cmdHelp(config *cliConfig, cmdArgs []string) error {
	fmt.Printf("\nWelcome to Gator, the Aggregator!\nUsage:\n\n")
	commands := getCommands()

	for _, command := range commands {
		fmt.Printf(" %v\t%v\n", command.name, command.description)
	}
	fmt.Println()
	return nil
}
