package cmd_test

import (
	"fmt"
	"github.com/kisielk/cmd"
	"os"
	"strings"
)

func Example_std() {
	hello := func(args []string) (string, error) {
		if len(args) == 0 {
			return "What's your name?\n", nil
		}
		return fmt.Sprintf("Hello, %s\n", strings.Join(args, " ")), nil
	}

	c := cmd.New(map[string]cmd.CmdFn{"hello": hello}, os.Stdin, os.Stdout)
	c.Loop()
}
