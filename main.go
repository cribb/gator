package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cribb/gator/internal/config"
	"github.com/cribb/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Something got jacked up while loading config.\n")
		os.Exit(1)
	}
	// fmt.Printf("cfg: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DbUrl)
	// fmt.Printf("db: %v\n", db)

	dbQueries := database.New(db)

	st := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	// fmt.Printf("st: %v\n", st)

	// Create the commands registry specified by `commands` type
	cmds := commands{
		handlers: map[string]func(*state, command) error{},
	}

	// Register the list of available commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	// Parse and prep the command entered on the CLI
	cmdEntered := os.Args
	if len(cmdEntered) < 2 {
		fmt.Println("No command provided.")
		os.Exit(1)
	}
	cmd := command{
		name: cmdEntered[1],
		args: cmdEntered[2:],
	}

	// Run the command and bounce
	if err := cmds.run(&st, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}
