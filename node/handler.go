package node

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
)

// TODO

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
	// TODO
	authed := G_SSHMgr.IsAuthorized(s.Conn().RemotePeer())
	if !authed {
		err = fmt.Errorf("handleStream: "+
			"RemotePeer not authorized.",
			s.Conn().RemotePeer().Pretty())
		return
	}

	return
}
