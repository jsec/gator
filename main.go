package main

import (
	"fmt"
	"log"

	"github.com/jsec/gator/internal/config"
	_ "github.com/lib/pq"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	// TODO: implement this
}

func (c *commands) run(s *state, cmd command) error {
	// TODO: implement this
	return nil
}

type state struct {
	config *config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("A name must be provided.")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error setting user:", err)
	}

	fmt.Println("User has been set.")
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	cfg.Print()

	err = cfg.SetUser("me")

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	cfg.Print()
}
