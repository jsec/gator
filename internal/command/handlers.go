package command

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jsec/gator/internal/database"
	"github.com/jsec/gator/internal/rss"
	"github.com/jsec/gator/internal/state"
)

func handlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("A name must be provided.")
	}

	user, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("User not found", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user:", err)
	}

	fmt.Println("User has been set.")
	return nil
}

func handlerRegister(s *state.State, cmd Command) error {
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

func handlerReset(s *state.State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting users:", err)
	}

	return nil
}

func handlerListUsers(s *state.State, cmd Command) error {
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

func handlerAggregate(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Not enough arguments provided")
	}

	rss.ScrapeFeeds(context.Background(), s.DB, cmd.Args[0])
	return nil
}

func handlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enough arguments provided")
	}

	ctx := context.Background()

	feed, err := s.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("Error creating feed:", err)
	}

	_, err = s.DB.CreateFollow(ctx, database.CreateFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("Error creating follow for feed:", err)
	}

	return nil
}

func handlerListAllFeeds(s *state.State, cmd Command) error {
	feeds, err := s.DB.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds:", err)
	}

	for _, feed := range feeds {
		fmt.Println("Feed:", feed.Name, "URL:", feed.Url, "User:", feed.UserName)
	}

	return nil
}

func handlerFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Not enough arguments provided")
	}

	ctx := context.Background()

	feed, err := s.DB.GetFeedByURL(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error fetching feed:", err)
	}

	follow, err := s.DB.CreateFollow(context.Background(), database.CreateFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	fmt.Println("Feed:", follow.FeedName, "User:", follow.UserName)
	return nil
}

func handlerGetUserFollows(s *state.State, cmd Command, user database.User) error {
	ctx := context.Background()

	follows, err := s.DB.GetFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("Error fetching follows for user:", err)
	}

	for _, follow := range follows {
		fmt.Println(follow.Name)
	}

	return nil
}

func handlerDeleteFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Not enough arguments provided")
	}

	err := s.DB.DeleteFollow(context.Background(), database.DeleteFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("Error deleting follow:", err)
	}

	return nil
}

func handlerBrowse(s *state.State, cmd Command, user database.User) error {
	var limit int32 = 2

	if len(cmd.Args) > 0 {
		argLimit, _ := strconv.Atoi(cmd.Args[0])
		limit = int32(argLimit)
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})

	if err != nil {
		return fmt.Errorf("Error fetching posts:", err)
	}

	for _, post := range posts {
		fmt.Println("Title:", post.Title, "URL:", post.Url)
	}

	return nil
}
