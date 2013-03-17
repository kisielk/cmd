// This is a simple example of using cmd for a command-line console.
package main

import (
	"fmt"
	"github.com/kisielk/cmd"
	"os"
	"strings"
)

func hello(args []string) (string, error) {
	if len(args) == 0 {
		return "What's your name?\n", nil
	}
	return fmt.Sprintf("Hello, %s\n", strings.Join(args, " ")), nil
}

func main() {
	c := cmd.New(map[string]cmd.CmdFn{"hello": hello}, os.Stdin, os.Stdout)
	c.Loop()
}
