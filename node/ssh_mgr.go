package node

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

// read in ~/.ipssh/authorized_keys

// TODO cache authorized keys

var G_SSHMgr *SSHMgr // TODO PROPERLY INIT

type SSHMgr struct {
	// TODO
}

func (s *SSHMgr) IsAuthorized(pid peer.ID) (tf bool) {
	// TODO
	return
}
