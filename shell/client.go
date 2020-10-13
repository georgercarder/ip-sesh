package shell

import (
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// TODO PUT IN COMMON
const StartShellSession = "Start Shell Session"

// trims byte(0), byte(0xa) byte(0x20)
func Trim(s string) (ret string) {
	// TODO PUT IN COMMON BUT IS HERE
	// SINCE IS USED TO COMPARE STARTSHELLSESSION
	lastByte := byte(s[len(s)-1])
	if lastByte == byte(0x0) || lastByte == byte(0xa) || lastByte == byte(0x20) {
		return Trim(s[:len(s)-1])
	}
	ret = s
	return

}

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
