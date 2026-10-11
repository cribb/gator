package main

import (
	// "bufio"
	// "encoding/json"
	// "strings"
	"context"
	"fmt"
	"time"

	"github.com/cribb/gator/internal/config"
	"github.com/cribb/gator/internal/database"
	"github.com/google/uuid"
	// "go.opentelemetry.io/otel/semconv/v1.43.0/dbconv"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	// name        string
	// description string
	handlers map[string]func(*state, command) error
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

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: no username provided.")
	}

	providedUserName := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), providedUserName)
	if err != nil {
		return fmt.Errorf("Error when logging in user '%v': (%w)", providedUserName, err)
	}

	// fmt.Printf("cmd.args: %v\n", cmd.args)
	// user := cmd.args[0]
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		// fmt.Println("failed to set user.")
		return err
	}
	fmt.Printf("User set to: %v\n", user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: no username provided.")
	} else if len(cmd.args) > 1 {
		return fmt.Errorf("Error: too many parameters.")
	}

	newUserName := cmd.args[0]
	now := time.Now().UTC()

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      newUserName,
	}
	// fmt.Printf("params: %v\n", params)

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error: failed when creating database user '%v'.\n  (%w)", user.Name, err)
	}
	fmt.Printf("created user '%v'.\n", user.Name)
	// fmt.Printf("userdata: %v\n", user)

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		// fmt.Println("failed to set user.")
		return err
	}
	fmt.Printf("User set to: %v\n", user.Name)
	return nil
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
