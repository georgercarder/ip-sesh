package node

import (
	"crypto/ed25519"
	//"encoding/json"
	"fmt"
	//"io"

	sh "github.com/georgercarder/ip-sesh/shell"

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
		/*case Shell:
		return shell(s)*/
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
		fmt.Println("debug err", err)
		return
	}
	fmt.Println("debug handshakePacket", hp)
	response, err := prepareChallengeResponse(hp)
	if err != nil {
		return
	}
	fmt.Println("debug before sendToStream", response)
	err = sendResponse(s, response.Bytes())
	fmt.Println("debug err", err)
	if err != nil {
		return
	}
	return checkHandshakeResult(s)
}

// 4
func checkHandshakeResult(s network.Stream) (err error) {
	fmt.Println("debug checkHandshakeResult start")
	ok, err := readHandshakeResult(s)
	fmt.Println("debug readHandshakeResult", ok, err)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("checkHandshakeResult: " +
			"Handshake result not ok.")
		return
	}
	fmt.Println("debug start shell session")
	return sh.Client(StreamToConn(s)) // shell session
}

// server

// 1
func sendChallenge(s network.Stream,
	nonce []byte, pk *ed25519.PublicKey) (err error) {
	chlg := prepareChallenge()
	hp := &HandshakePacket{Challenge: chlg, PubKey: pk, Nonce: nonce}
	G_HandshakeMgr.newHandshake(hp, pk)
	err = sendToStream(s, HandshakeInitChallenge, hp.Bytes())
	if err != nil {
		return
	}
	return checkHandshakeResponse(s)
}

// 3
func checkHandshakeResponse(s network.Stream) (err error) {
	fmt.Println("debug checkHandshakeResponse")
	hp, err := readHandshakePacket(s)
	if err != nil {
		fmt.Println("debug err", err)
		return
	}
	fmt.Println("debug hp", hp)
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
	err = sh.Server(StreamToConn(s))
	return
}
