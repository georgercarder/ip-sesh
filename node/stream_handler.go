package node

import (
	"crypto/ed25519"
	//"encoding/json"
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
)

func StreamHandler(s network.Stream) {
	//fmt.Println("debug new stream!")
	if err := handleStream(s); err != nil {
		//LogError.Println("StreamHandler:", err)
		s.Reset()
	} else {
		s.Close()
	}
}

func handleStream(s network.Stream) (err error) {
	ss := make([]byte, 1)
	_, err = s.Read(ss)
	if err != nil {
		return
	}
	switch StreamStatus(ss[0]) {
	case HandshakeInitChallenge:
		return handleHandshakeInitChallenge(s)
	case HandshakeResponse:
		return checkHandshakeResponse(s)
	case HandshakeResult:
		return checkHandshakeResult(s)
	case ShellFrame:
		return shellFrame(s)
	}
	return fmt.Errorf("handleStream: "+
		"StreamStatus not supported. %d", ss[0])
}

// client

// 2
func handleHandshakeInitChallenge(s network.Stream) (err error) {
	fmt.Println("debug handleHandshakeInitChallenge")
	hp, err := readHandshakePacket(s)
	if err != nil {
		return
	}
	response, err := prepareChallengeResponse(hp)
	if err != nil {
		return
	}
	return sendToStream(s, HandshakeResponse, response.Bytes())
}

// 4
func checkHandshakeResult(s network.Stream) (err error) {
	ok, err := readHandshakeResult(s)
	if err != nil {
		return
	}
	if ok != true {
		err = fmt.Errorf("checkHandshakeResult: " +
			"Handshake result not ok.")
	}
	// TODO RUN SHELL SESSION
	fmt.Println("debug run shell session here")
	return
}

// server

// 1
func sendChallenge(s network.Stream, pk *ed25519.PublicKey) (err error) {
	chlg := prepareChallenge()
	hp := &HandshakePacket{Challenge: chlg, PubKey: pk}
	return sendToStream(s, HandshakeInitChallenge, hp.Bytes())
}

// 3
func checkHandshakeResponse(s network.Stream) (err error) {
	hp, err := readHandshakePacket(s)
	if err != nil {
		return
	}
	ok := checkAgainstPendingHandshakes(hp)
	g_activeSessions.Put(s.Conn().RemotePeer())
	return sendBackToClient(s, HandshakeResult,
		[]byte{Boole2Byte(ok)})
}

func shellFrame(s network.Stream) (err error) {
	ok := g_activeSessions.Check(s.Conn().RemotePeer())
	if !ok {
		err = fmt.Errorf("shellFrame: "+
			"not an active session. %q",
			s.Conn().RemotePeer())
		return sendBackToClient(s, Error, nil)
	}
	return runSession(s)
}

func runSession(s network.Stream) (err error) {
	sh := initShell()
	go func() {
		// send stream back stdout, stderr
		_, er := io.Copy(s, sh)
		if er != nil {
			if er != io.EOF {
				// TODO LOG
				return
			}
		}
	}()
	_, err = io.Copy(sh, s)
	return
}
