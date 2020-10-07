package node

import (
	"crypto/ed25519"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type StreamStatus byte

const (
	HandshakeInitChallenge StreamStatus = iota
	HandshakeResponse
	HandshakeResult
	ShellFrame
	Error
)

func (ss StreamStatus) Byte() (b byte) {
	b = byte(ss)
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

func (p *pendingHandshakes) Put(b []byte, pk *ed25519.PublicKey) {
	// TODO
	return
}

func (p *pendingHandshakes) Check(hp *HandshakePacket) (ok bool) {
	// TODO
	// checks if ed25519 sig is valid

	return
}

func checkAgainstPendingHandshakes(hp *HandshakePacket) (ok bool) {
	return g_pendingHandshakes.Check(hp)
}

func sendBackToClient(
	s network.Stream, status StreamStatus, msg []byte) (err error) {
	return writeToStream(s, status, msg)
}

func sendToStream(
	s network.Stream, status StreamStatus, msg []byte) (err error) {
	return writeToStream(s, status, msg)
}

func writeToStream(
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

func readHandshakePacket(s network.Stream) (hp *HandshakePacket, err error) {
	// TODO
	return
}

func readHandshakeResult(s network.Stream) (res bool, err error) {
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
