package main

import (
	"fmt"
	"os"

	"github.com/cribb/gator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Something got jacked up while loading config.\n")
		os.Exit(1)
	}
	// fmt.Printf("cfg: %+v\n", cfg)

	st := state{
		cfg: &cfg,
	}
	// fmt.Printf("st: %v\n", st)

	// commands := map[string]cliCommand
	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)

	cmdEntered := os.Args
	if len(cmdEntered) < 2 {
		fmt.Println("No command provided.")
		os.Exit(1)
	}

	cmd := command{
		name: cmdEntered[1],
		args: cmdEntered[2:],
	}

	if err := cmds.run(&st, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}
