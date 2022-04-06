package commands

import (
	"fmt"
	"log"

	"github.com/rafaelfuchsbr/tcp-service/internal/globals"
)

type CommandWho struct {
}

func NewCommandWho() *CommandWho {
	return &CommandWho{}
}

func (c *CommandWho) Process() (string, error) {
	response := fmt.Sprintf("%d", globals.ConnectedClients.Total())
	log.Println(fmt.Sprintf("Response: %s", response))
	return response, nil
}
