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
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
	}()
	_, err = io.Copy(conn, os.Stdin)
	return err
}
