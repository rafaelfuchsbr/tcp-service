package commands

import (
	"fmt"
	"log"
	"strings"
)

const (
	COMMAND_WHERE = "WHERE"
	COMMAND_WHO   = "WHO"
	COMMAND_WHY   = "WHY"
)

type Processor interface {
	Process() (string, error)
}

func NewCommand(input string) (Processor, error) {

	input = strings.Trim(input, " \n\r")

	log.Println(fmt.Sprintf("Got command [%s] - generating new command object", input))

	switch input {
	case COMMAND_WHERE:
		return NewCommandWhere(), nil
	case COMMAND_WHO:
		return NewCommandWho(), nil
	case COMMAND_WHY:
		return NewCommandWhy(), nil
	}

	return nil, fmt.Errorf("Input is not a valid command: [%s]", input)
}
