package commands

import (
	"fmt"
	"log"
)

type CommandWhy struct {
}

func NewCommandWhy() *CommandWhy {
	return &CommandWhy{}
}

func (c *CommandWhy) Process() (string, error) {
	response := "42"
	log.Println(fmt.Sprintf("Response: %s", response))
	return response, nil
}
