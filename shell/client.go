package shell

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const StartShellSession = "Start Shell Session"

func Client(conn net.Conn) (err error) {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		fmt.Println("debug after1")
	}()
	_, err = io.Copy(conn, os.Stdin)
	fmt.Println("debug after2")
	return err
}
