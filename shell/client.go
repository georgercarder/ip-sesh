package shell

import (
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func Client(conn net.Conn) (err error) {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}() // Best effort.
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
	}()
	_, err = io.Copy(conn, os.Stdin)
	return err
}
