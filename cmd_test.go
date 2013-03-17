package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/kisielk/cmd"
	"testing"
)

func TestOne(t *testing.T) {
	out := &bytes.Buffer{}
	commands := map[string]cmd.CmdFn{
		"good": func(args []string) (string, error) {
			return fmt.Sprintf("good %v\n", args), nil
		},
		"bad": func(args []string) (string, error) {
			return "bad\n", fmt.Errorf("oops")
		},
	}
	c := cmd.New(commands, nil, out)

	tests := []struct {
		In          string
		Out         string
		ShouldError bool
	}{
		{"hello ", "unrecognized input: hello \n", false},
		{"  ", "unrecognized input: hello \n", false},
		{"good", "good []\n", false},
		{"good arg1 arg2", "good [arg1 arg2]\n", false},
		{"bad", "bad\n", true},
	}

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
