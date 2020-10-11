package shell

import (
	"io"
	"net"
	"os/exec"
	"os/user"

	"github.com/creack/pty"
)

func Server(conn net.Conn) error {
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
		return e
	}
	// Make sure to close the pty at the end.
	defer func() {
		_ = ptmx.Close()
	}() // Best effort.
	go func() {
		_, _ = io.Copy(ptmx, conn)
	}()
	_, e = io.Copy(conn, ptmx)
	return e
}
