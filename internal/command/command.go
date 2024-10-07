package command

import (
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

func (c *Commands) register(name string, f func(*state.State, Command) error) {
	c.Handlers[name] = f
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	// TODO: handle invalid commands
	handler := c.Handlers[cmd.Name]

	err := handler(s, cmd)
	if err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func NewCommands() Commands {
	commands := Commands{
		Handlers: make(map[string]func(*state.State, Command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAggregate)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerListAllFeeds)
	commands.register("follow", handlerFollow)
	commands.register("following", handlerGetUserFollows)

	return commands
}
