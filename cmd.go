// Package cmd provides a simple way to creates simple line-oriented interactive command interpreters.
// It's inspired by the cmd module from the Python standard library.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// DefaultPrompt is the default value of Cmd.Prompt
const DefaultPrompt = "> "

// CmdFn is the function type that can be used to define commands for a Cmd.
// The value of out is printed to the console.
// If err is not nil then execution of the command loop is terminated.
type CmdFn func(args []string) (out string, err error)

// Cmd is an interactive command interpreter. It's started by calling the Loop method.
// Instances of Cmd should be constructed with the New function.
type Cmd struct {
	// In receives input
	In io.Reader

	// Out transmits output
	Out io.Writer

	// Prompt is displayed on the console before every line of input.
	Prompt string

	// Commands is a map of command functions for valid commands.
	// If a command is not in this map then Default will be called.
	Commands map[string]CmdFn

	// Default is called when a command is received that does not match
	// any function in the Commands map. The line argument will contain
	// the full contents of the line received.
	//
	// The value of out is printed to the console.
	// If err is not nil then execution of the command loop is terminated.
	//
	// If Default is not set the behaviour is to print a message to the
	// console.
	Default func(line string) (out string, err error)

	// EmptyLine is called whenever a line containing no characters
	// other than whitespace or newline is received.
	//
	// The value of out is printed to the console.
	// If err is not nil then execution of the command loop is terminated.
	//
	// If EmptyLine is not set then the last command is repeated.
	EmptyLine func() (out string, err error)

	// Tokens is called for each line of input to generate the tokens.
	//
	// The first token is the name of the command that will be called,
	// while the rest of the tokens are passed as arguments to the command.
	//
	// If Tokens is not set then strings.Fields is used.
	Tokens func(line string) (tokens []string)

	// LastLine contains the last non-empty line received
	LastLine string
}

// New creates a new Cmd with the commands from c that communicates via in and out.
func New(c map[string]CmdFn, in io.Reader, out io.Writer) *Cmd {
	cmd := Cmd{In: in, Out: out, Prompt: DefaultPrompt, LastLine: "", Commands: c}
	cmd.EmptyLine = func() (string, error) {
		if len(cmd.LastLine) > 0 {
			return "", cmd.one(cmd.LastLine)
		}
		return "", nil
	}
	cmd.Default = func(line string) (string, error) {
		return fmt.Sprintf("unrecognized command: %s\n", strings.Fields(line)[0]), nil
	}
	cmd.Tokens = strings.Fields
	return &cmd
}

func (c *Cmd) parseLine(line string) (cmd string, args []string) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return
	}

	tokens := c.Tokens(line)
	if len(tokens) == 0 {
		return
	}
	cmd = tokens[0]
	if len(tokens) > 1 {
		args = tokens[1:]
	}
	return
}

// one parses one line of input and executes a command.
// The output of the command is sent to c.Out.
func (c *Cmd) one(line string) error {
	cmd, args := c.parseLine(line)

	var msg string
	var cmderr error

	if cmd == "" {
		msg, cmderr = c.EmptyLine()
	} else {
		c.LastLine = line[:]
		if fn := c.Commands[cmd]; fn == nil {
			msg, cmderr = c.Default(line)
		} else {
			msg, cmderr = fn(args)
		}
	}

	if msg != "" {
		if _, err := c.Out.Write([]byte(msg)); err != nil {
			return err
		}
	}
	return cmderr
}

// Loop starts the interpreter loop.
//
// For each iteration it prints c.Prompt to c.Out and then waits for a line of input.
// The line is tokenized using c.Tokens and the first token is interpreted as the name
// of a command. The command is looked up in c.Commands and is called with the remaining
// tokens.
// 
// If the command is not found then c.Default is called with the entire line.
//
// If the input line consists only of whitespace then c.EmptyLine is called.
func (c *Cmd) Loop() error {
	rd := bufio.NewReader(c.In)
	for {
		_, err := c.Out.Write([]byte(c.Prompt))
		if err != nil {
			return err
		}
		line, err := rd.ReadBytes('\n')
		if err != nil {
			return err
		}
		if err := c.one(string(line)); err != nil {
			return err
		}
	}
	panic("unreachable")
}
