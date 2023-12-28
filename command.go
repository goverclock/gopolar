package gopolar

import (
	"strings"
)

type Command struct {
	Name string
	args []string
}

func NewCommand(s string) *Command {
	split := strings.Fields(s)
	return &Command{
		Name: split[0],
		args: split[1:],
	}
}
