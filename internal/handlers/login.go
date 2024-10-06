package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jsec/gator/internal/command"
	"github.com/jsec/gator/internal/database"
	"github.com/jsec/gator/internal/state"
)

func Login(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("A name must be provided.")
	}

	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("That user does not exist", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user:", err)
	}

	fmt.Println("User has been set.")
	return nil
}

func Register(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("A name must be provided.")
	}

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("Error creating user:", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting current user:", err)
	}

	return nil
}
