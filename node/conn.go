package node

import (
	"net"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
)

type Conn struct {
	S network.Stream
}

func StreamToConn(s network.Stream) (c *Conn) {
	c = new(Conn)
	c.S = s
	return
}

func (c *Conn) Read(b []byte) (int, error) {
	return c.S.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.S.Write(b)
}

func (c Conn) Close() (err error) {
	return c.S.Close()
}

func (c Conn) RemoteAddr() (addr net.Addr) {
	// TODO
	return
}

func (c Conn) SetDeadline(t time.Time) (err error) {
	// TODO
	return
}

func (c Conn) SetReadDeadline(t time.Time) (err error) {
	// TODO
	return
}

func (c Conn) SetWriteDeadline(t time.Time) (err error) {
	// TODO
	return
}

func (c *Conn) LocalAddr() (addr net.Addr) {
	// TODO
	return
}
