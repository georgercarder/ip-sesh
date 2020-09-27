package node

import (
	"fmt"
	"io"

	"github.com/libp2p/go-libp2p-core/network"
)

// TODO in progress

// will handle streams

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
	authed := G_SSHMgr.IsAuthorized(s.Conn().RemotePeer())
	if !authed {
		err = fmt.Errorf("handleStream: "+
			"RemotePeer not authorized.",
			s.Conn().RemotePeer().Pretty())
		return
	}
	return startSession(s)
}

func startSession(s network.Stream) (err error) {
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

func initShell() (s *shell) {
	// TODO
	return
}

type shell struct {
	// TODO
}

func (s *shell) Read(b []byte) (n int, err error) {
	// TODO
	return
}

func (s *shell) Write(b []byte) (n int, err error) {
	// TODO
	return
}
