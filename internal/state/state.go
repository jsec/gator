package state

import (
	"log"

	"github.com/jsec/gator/internal/config"
	"github.com/jsec/gator/internal/database"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

func New() State {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	return State{
		Config: &cfg,
	}
}
