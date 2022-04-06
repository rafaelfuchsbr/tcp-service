package commands

import (
	"fmt"
	"log"

	"github.com/rafaelfuchsbr/tcp-service/internal/globals"
)

type CommandWhere struct {
}

func NewCommandWhere() *CommandWhere {
	return &CommandWhere{}
}

func (c *CommandWhere) Process() (string, error) {
	response := globals.ServerId
	log.Println(fmt.Sprintf("Response: %s", response))
	return response, nil
}
