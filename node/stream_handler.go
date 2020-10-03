package node

import (
	"crypto/ed25519"
	"encoding/json"
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
	/*case InitStream:
	return initStream(s)*/
	case HandshakeInitPublicKeys:
		// TODO
		return nil
	case HandshakeInitChallenge:
		// TODO
		return nil
	case HandshakeResponse:
		return checkHandshakeResponse(s)
	case ShellFrame:
		return shellFrame(s)
	}
	return fmt.Errorf("handleStream: "+
		"StreamStatus not supported.", ss[0])
}

// TODO DELETE but currently keeping for reference during dev
/*func initStream(s network.Stream) (err error) {
	pk, err := readIPSSHPubKey(s)
	if err != nil {
		return
	}
	// note: this pk is to ipssh keyset which is
	// different from ipfs node credentials.
	// See main doc.
	if !G_SSHMgr.IsAuthorized(pk) {
		err = fmt.Errorf("handleStream: "+
			"RemotePeer not authorized.",
			s.Conn().RemotePeer().Pretty())
		// LogError.Println(err)
		return
	}
	return sendChallenge(s, pk)
}*/

func sendPublicKeys(s network.Stream, pks []*ed25519.PublicKey) (err error) {
	b, err := json.Marshal(pks)
	if err != nil {
		// TODO log
	}
	return sendBackToClient(s, HandshakeInitPublicKeys, b)
}

func sendChallenge(s network.Stream, pk *ed25519.PublicKey) (err error) {
	chlg, err := prepareChallenge()
	if err != nil {
		return
	}
	//g_pendingHandshakes.Put(chlg, pk)
	return sendBackToClient(s, HandshakeInitChallenge, chlg)
}

func checkHandshakeResponse(s network.Stream) (err error) {
	hsRes, err := readHandshakeResponse(s)
	if err != nil {
		return
	}
	ok := checkAgainstPendingHandshakes(hsRes)
	g_activeSessions.Put(s.Conn().RemotePeer())
	return sendBackToClient(s, HandshakeResponse,
		[]byte{Boole2Byte(ok)})
}

func shellFrame(s network.Stream) (err error) {
	ok := g_activeSessions.Check(s.Conn().RemotePeer())
	if !ok {
		err = fmt.Errorf("shellFrame: "+
			"not an active session.",
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
