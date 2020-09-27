package node

import (
	"crypto/ed25519"
)

// read in ~/.ipssh/authorized_keys

// TODO cache authorized keys

var G_SSHMgr *SSHMgr // TODO PROPERLY INIT

type SSHMgr struct {
	// TODO
}

func (s *SSHMgr) IsAuthorized(pk *ed25519.PublicKey) (tf bool) {
	// TODO
	return
}
