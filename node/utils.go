package node

import (
	"fmt"
	"io"

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

func checkAgainstPendingHandshakes(hp *HandshakePacket) (ok bool) {
	storedHs := G_HandshakeMgr.DomainName2Handshake.Take(hp.DomainName)
	if storedHs == nil {
		// LOG
		return
	}
	sHs := storedHs.(*Handshake)
	return VerifySig(sHs.PubKey, sHs.Challenge, hp.Signature)
}

func sendResponse(s network.Stream, msg []byte) (err error) {
	_, err = io.Copy(s, newMsgReader(msg))
	return
}

type msgReader struct {
	buf []byte
}

func newMsgReader(b []byte) io.Reader {
	m := &msgReader{buf: b}
	return m
}

func (m *msgReader) Read(b []byte) (n int, err error) {
	copy(b, m.buf[:])
	if len(b) > len(m.buf) {
		n = len(m.buf)
	} else {
		n = len(b)
	}
	if n == 0 {
		err = io.EOF
		return
	}
	m.buf = m.buf[n:]
	return
}

func sendBackToClient(
	s network.Stream, status StreamStatus, msg []byte) (err error) {
	return writeToStream(s, status, msg)
}

func sendToStream(
	s network.Stream, status StreamStatus, msg []byte) (err error) {
	fmt.Println("debug sendToStream", status, msg)
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

func Byte2Boole(b byte) (ret bool) {
	if b == byte(0x1) {
		ret = true
	}
	return
}

func readHandshakePacket(s network.Stream) (hp *HandshakePacket, err error) {
	b := make([]byte, 1024)
	n, err := s.Read(b)
	if err != nil {
		if err != io.EOF {
			return
		}
	}
	return Slice2HandshakePacket(b[:n])
}

func readHandshakeResult(s network.Stream) (res bool, err error) {
	b := make([]byte, 1)
	n, err := s.Read(b)
	if err != nil {
		if err != io.EOF {
			return
		}
	}
	if n < 1 {
		err = fmt.Errorf("Must read at least 1 byte.")
		return
	}
	res = Byte2Boole(b[0])
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
