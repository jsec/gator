package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jsec/gator/internal/command"
	"github.com/jsec/gator/internal/database"
	"github.com/jsec/gator/internal/state"
	_ "github.com/lib/pq"
)

func main() {
	s := state.New()

	db, err := sql.Open("postgres", s.Config.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database:", err)
	}
	defer db.Close()

	s.DB = database.New(db)
	commands := command.NewCommands()
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Not enough arguments specified")
	}

	commands.Run(&s, command.Command{
		Name: args[1],
		Args: args[2:],
	})
}
