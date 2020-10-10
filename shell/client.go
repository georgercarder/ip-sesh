package shell

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func Client() error {
	// connect to this socket
	conn, e := net.Dial("tcp", "127.0.0.1:8081")
	if e != nil {
		return e
	}

	// MakeRaw put the terminal connected to the given file descriptor into raw
	// mode and returns the previous state of the terminal so that it can be
	// restored.
	oldState, e := terminal.MakeRaw(int(os.Stdin.Fd()))
	if e != nil {
		return e
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() { _, _ = io.Copy(os.Stdout, conn) }()
	_, e = io.Copy(conn, os.Stdin)
	fmt.Println("Bye!")

	return e
}
