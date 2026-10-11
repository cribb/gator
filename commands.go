package main

import (
	// "bufio"
	// "encoding/json"
	// "strings"
	"fmt"

	"github.com/cribb/gator/internal/config"
)

type command struct {
	name string
	args []string
}

type commands struct {
	// name        string
	// description string
	handlers map[string]func(*state, command) error
}

type state struct {
	cfg *config.Config
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: no username provided.")
	}

	// fmt.Printf("cmd.args: %v\n", cmd.args)
	user := cmd.args[0]
	err := s.cfg.SetUser(user)
	if err != nil {
		// fmt.Println("failed to set user.")
		return err
	}
	fmt.Printf("User set to: %v\n", user)
	return nil
}

func (c *commands) run(s *state, cmd command) error {

	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s not found.", cmd.name)
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// func (c *commands) handlerHelp(config *cliConfig, cmdArgs []string) error {
// 	fmt.Printf("\nWelcome to Gator, the Aggregator!\nUsage:\n\n")
// 	commands := getCommands()

// 	for _, command := range commands {
// 		fmt.Printf(" %v\t%v\n", command.name, command.description)
// 	}
// 	fmt.Println()
// 	return nil
// }
