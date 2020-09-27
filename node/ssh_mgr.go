package node

import (
	"crypto/rsa"
)

// read in ~/.ipssh/authorized_keys

// TODO cache authorized keys

var G_SSHMgr *SSHMgr // TODO PROPERLY INIT

type SSHMgr struct {
	// TODO
}

func (s *SSHMgr) IsAuthorized(pk *rsa.PublicKey) (tf bool) {
	// TODO
	return
}
