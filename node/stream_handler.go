package node

import (
	"crypto/ed25519"
	"fmt"

	sh "github.com/georgercarder/ip-sesh/shell"

	"github.com/libp2p/go-libp2p-core/network"
)

func StreamHandler(s network.Stream) {
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
	}
	return fmt.Errorf("handleStream: "+
		"StreamStatus not supported. %d", ss[0])
}

// client

// 2
func handleHandshakeInitChallenge(s network.Stream) (err error) {
	hp, err := readHandshakePacket(s)
	if err != nil {
		fmt.Println("debug err", err)
		return
	}
	response, err := prepareChallengeResponse(hp)
	if err != nil {
		return
	}
	err = sendResponse(s, response.Bytes())
	if err != nil {
		return
	}
	return checkHandshakeResult(s, hp)
}

// 4
func checkHandshakeResult(s network.Stream, hp *HandshakePacket) (err error) {
	ok, err := readHandshakeResult(s)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("checkHandshakeResult: " +
			"Handshake result not ok.")
		return
	}
	// TODO PASS STREAMPACKETCH HERE
	h := G_HandshakeMgr.Nonce2Handshake.Take(string(hp.Nonce))
	if h == nil {
		err = fmt.Errorf("nonce does not map to a handshake.")
		return
	}
	hs := h.(*Handshake)
	stopCH := make(chan bool, 1)
	hs.ConnBundleCH <- &ConnBundle{Conn: StreamToConn(s), StopCH: stopCH}
	_ = <-stopCH // holds stream open while clientDaemon uses it
	fmt.Println("debug stopCH")
	//return sh.Client(StreamToConn(s)) // shell session
	return
}

// server

// 1
func sendChallenge(s network.Stream,
	nonce []byte, pk *ed25519.PublicKey) (err error) {
	chlg := prepareChallenge()
	hp := &HandshakePacket{Challenge: chlg, PubKey: pk, Nonce: nonce}
	G_HandshakeMgr.newHandshake(hp, pk, nil)
	err = sendToStream(s, HandshakeInitChallenge, hp.Bytes())
	if err != nil {
		return
	}
	return checkHandshakeResponse(s)
}

// 3
func checkHandshakeResponse(s network.Stream) (err error) {
	hp, err := readHandshakePacket(s)
	if err != nil {
		fmt.Println("debug err", err)
		return
	}
	ok := checkAgainstPendingHandshakes(hp)
	g_activeSessions.Put(s.Conn().RemotePeer())
	err = sendResponse(s, []byte{Boole2Byte(ok)})
	if err != nil {
		fmt.Println("debug err", err)
		return
	}
	return shell(s)
}

func shell(s network.Stream) (err error) {
	ok := g_activeSessions.Check(s.Conn().RemotePeer())
	if !ok {
		err = fmt.Errorf("shell: "+
			"not a valid active session. %q",
			s.Conn().RemotePeer())
		return sendBackToClient(s, Error, nil)
	}
	return runSession(s)
}

func runSession(s network.Stream) (err error) {
	return sh.Server(StreamToConn(s))
}
