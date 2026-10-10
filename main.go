package main

import (
	"fmt"

	"github.com/cribb/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	cfg.SetUser("jeremy")
	if err != nil {
		fmt.Printf("Something got jacked up while loading config.\n")
		return
	}
	// cfg, err = config.Read()
	// fmt.Printf("cfg: %+v\n", cfg)
}
