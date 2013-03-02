// This is a simple example of using cmd.go over a tcp socket.
// Start the example and then telnet to the host on port 6000 to see it in action.
package main

import (
	"fmt"
	"github.com/kisielk/cmd.go"
	"log"
	"net"
	"strings"
)

func hello(args []string) (string, error) {
	if len(args) == 0 {
		return "What's your name?\n", nil
	}
	return fmt.Sprintf("Hello, %s\n", strings.Join(args, " ")), nil
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal("could not open port:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("couldn't accept console:", err)
			continue
		}
		c := cmd.New(map[string]cmd.CmdFn{"hello": hello}, conn, conn)
		go c.Loop()
	}
}
