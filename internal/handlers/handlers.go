package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jsec/gator/internal/command"
	"github.com/jsec/gator/internal/database"
	"github.com/jsec/gator/internal/rss"
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

func Reset(s *state.State, cmd command.Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting users:", err)
	}

	return nil
}

func ListUsers(s *state.State, cmd command.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving users:", err)
	}

	for _, user := range users {
		name := user.Name

		if name == s.Config.CurrentUserName {
			name += " (current)"
		}

		fmt.Println("*", name)
	}

	return nil
}

func AggregateFeeds(s *state.State, cmd command.Command) error {
	feed, err := rss.GetFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error retrieving RSS feed:", err)
	}

	fmt.Println(feed)
	return nil
}

func AddFeed(s *state.State, cmd command.Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough arguments provided")
	}

	ctx := context.Background()

	user, err := s.DB.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error fetching current user:", err)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	fmt.Println(feed)
	return nil
}

func ListAllFeeds(s *state.State, cmd command.Command) error {
	feeds, err := s.DB.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds:", err)
	}

	for _, feed := range feeds {
		fmt.Println("Feed:", feed.Name, "URL:", feed.Url, "User:", feed.UserName)
	}

	return nil
}
