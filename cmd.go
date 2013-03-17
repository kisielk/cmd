// Package cmd provides a simple way to creates simple line-oriented interactive command interpreters.
// It's inspired by the cmd module from the Python standard library.
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// DefaultPrompt is the default value of Cmd.Prompt
const DefaultPrompt = "> "

// CmdFn is the function type that can be used to define commands for a Cmd.
// The value of out is printed to the console.
// If err is not nil then execution of the command loop is terminated.
type CmdFn func(args []string) (out string, err error)

// Cmd is an interactive command interpreter. It's started by calling the Loop method.
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
	Default func(line []byte) (out string, err error)

	// EmptyLine is called whenever a line containing no characters
	// other than whitespace or newline is received.
	//
	// The value of out is printed to the console.
	// If err is not nil then execution of the command loop is terminated.
	//
	// If EmptyLine is not set then the last command is repeated.
	EmptyLine func() (out string, err error)

	// LastLine contains the last non-empty line received
	LastLine []byte
}

// New creates a new Cmd with the commands from c that communicates via in and out.
func New(c map[string]CmdFn, in io.Reader, out io.Writer) *Cmd {
	cmd := Cmd{In: in, Out: out, Prompt: DefaultPrompt, LastLine: make([]byte, 0), Commands: c}
	cmd.EmptyLine = func() (string, error) {
		if len(cmd.LastLine) > 0 {
			return "", cmd.One(cmd.LastLine)
		}
		return "", nil
	}
	cmd.Default = func(line []byte) (string, error) {
		return fmt.Sprintf("unrecognized input: %s\n", line), nil
	}
	return &cmd
}

func (c *Cmd) parseLine(line []byte) (cmd string, args []string) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return
	}

	fields := bytes.Fields(line)
	if len(fields) == 0 {
		return
	}
	cmd = string(fields[0])
	for _, f := range fields[1:] {
		args = append(args, string(f))
	}
	return
}

// One parses one line of input and executes a command.
// The output of the command is sent to c.Out.
func (c *Cmd) One(line []byte) error {
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
// For each iteration it prints c.Prompt to c.Out and then waits for a line of input.
// The line is used to call c.One. The loop terminates and returns an error if
// any call to c.One fails.
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
		if err := c.One(line); err != nil {
			return err
		}
	}
	panic("unreachable")
}
