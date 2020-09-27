package node

import (
	"crypto/rsa"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type StreamStatus byte

const (
	InitStream StreamStatus = iota
	HandshakeResponse
	ShellFrame
	Error
)

func (ss StreamStatus) Byte() (b byte) {
	b = byte(ss)
	return
}

func readIPSSHPubKey(s network.Stream) (pk *rsa.PublicKey, err error) {
	// TODO
	return
}

var g_pendingHandshakes = newPendingHandshakes()

func newPendingHandshakes() (p *pendingHandshakes) {
	// TODO init
	return
}

type pendingHandshakes struct {
	// TODO
}

func (p *pendingHandshakes) Put(b []byte, pk *rsa.PublicKey) {
	// TODO
	return
}

func (p *pendingHandshakes) Check(b []byte) (ok bool) {
	// TODO
	return
}

func checkAgainstPendingHandshakes(hsRes []byte) (ok bool) {
	return g_pendingHandshakes.Check(hsRes)
}

func sendBackToClient(
	s network.Stream, status StreamStatus, msg []byte) (err error) {
	fullMsg := append([]byte{status.Byte()}, msg...)
	_, err = s.Write(fullMsg)
	return err
}

func Boole2Byte(b bool) (ret byte) {
	if b {
		ret = byte(1)
	}
	return
}

func readHandshakeResponse(s network.Stream) (res []byte, err error) {
	// TODO
	return
}

var g_activeSessions = newActiveSessions()

func newActiveSessions() (p *activeSessions) {
	// TODO init
	return
}

type activeSessions struct {
	// TODO
}

func (p *activeSessions) Put(pid peer.ID) {
	// TODO
	return
}

func (p *activeSessions) Check(pid peer.ID) (ok bool) {
	// TODO
	return
}
