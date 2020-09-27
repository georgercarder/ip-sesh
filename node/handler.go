package node

import (
	"github.com/libp2p/go-libp2p-core/network"
)

// TODO

// will handle streams

// check attempt at connection against
// authorized_keys from ~/.ipssh/authorized_keys

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
	return
}

