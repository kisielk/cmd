package cmd_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kisielk/cmd"
	"io"
	"testing"
)

var commands = map[string]cmd.CmdFn{
	"good": func(args []string) (string, error) {
		return fmt.Sprintf("good %v\n", args), nil
	},
	"bad": func(args []string) (string, error) {
		return "bad\n", fmt.Errorf("oops")
	},
}

var tests = []struct {
	In          string
	Out         string
	ShouldError bool
}{
	{"hello ", "unrecognized command: hello\n", false},
	{"  ", "unrecognized command: hello\n", false},
	{"good", "good []\n", false},
	{"good arg1 arg2", "good [arg1 arg2]\n", false},
	{"bad", "bad\n", true},
}

func TestOne(t *testing.T) {
	out := &bytes.Buffer{}
	c := cmd.New(commands, nil, out)

	for i, test := range tests {
		if err := c.One([]byte(test.In)); !test.ShouldError && err != nil {
			t.Fatalf("%d: unexpected error: %s", i, err)
		} else if test.ShouldError && err == nil {
			t.Fatalf("%d: expected error but got nil")
		}

		if outMsg := out.String(); outMsg != test.Out {
			t.Fatalf("%d: bad output: got %q, want %q", i, outMsg, test.Out)
		}
		out.Reset()
	}
}

func TestLoop(t *testing.T) {
	in, inw := io.Pipe()
	outr, out := io.Pipe()
	outbuf := bufio.NewReader(outr)
	c := cmd.New(commands, in, out)

	go func() {
		err := c.Loop()
		if err != nil {
			t.Fatal(err)
		}
	}()

	promptbuf := make([]byte, len(c.Prompt))
	for i, test := range tests {
		if test.ShouldError {
			// Don't run tests that would stop the loop
			continue
		}

		_, err := outbuf.Read(promptbuf)
		if err != nil {
			t.Fatalf("%d: couldn't read prompt: %s", i, err)
		}
		if prompt := string(promptbuf); prompt != c.Prompt {
			t.Fatalf("%d: bad prompt, got %q want %q", i, prompt, c.Prompt)
		}
		fmt.Fprintln(inw, test.In)
		outMsg, err := outbuf.ReadBytes('\n')
		if err != nil {
			t.Fatalf("%d: couldn't read output: %s", err)
		}
		if o := string(outMsg); o != test.Out {
			t.Fatalf("%d: bad output: got %q, want %q", i, o, test.Out)
		}
	}
}
