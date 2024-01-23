package bot

import "fmt"

type command struct {
	handle func(args []string) error
	hint   string
}

const PREFIX = "!"

var actions = map[string]command{
	prepend("help"): {
		handle: func(args []string) error {
			fmt.Println(args)

			return nil
		},
		hint: "Print a help message",
	},
}

func prepend(cmd string) string {
	return fmt.Sprintf("%s%s", PREFIX, cmd)
}
