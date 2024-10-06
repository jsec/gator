package command

import (
	"fmt"
	"log"

	"github.com/jsec/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*state.State, Command) error
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.Handlers[name] = f
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	// TODO: handle invalid commands
	handler := c.Handlers[cmd.Name]

	err := handler(s, cmd)
	if err != nil {
		fmt.Println("Error:", err.Error())
		log.Fatal(err.Error())
	}

	return nil
}
