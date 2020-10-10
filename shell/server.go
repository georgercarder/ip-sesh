package shell

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"

	"github.com/creack/pty"
)

func Server() error {
	// Create command
	c := exec.Command("bash") // TODO put in config preferred shell
	usr, err := user.Current()
	if err != nil {
		return err
	}
	c.Dir = usr.HomeDir
	// Start the command with a pty.
	ptmx, e := pty.Start(c)
	if e != nil {
		fmt.Println("debug err", e)
		return e
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	return listen(ptmx)
}

func listen(ptmx *os.File) error {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, e := net.Listen("tcp", ":8081")
	if e != nil {
		return e
	}
	// accept connection on port
	conn, e := ln.Accept()
	if e != nil {
		return e
	}
	fmt.Println("debug accept")
	go func() {
		_, _ = io.Copy(ptmx, conn)
	}()
	_, e = io.Copy(conn, ptmx)
	return e
}
