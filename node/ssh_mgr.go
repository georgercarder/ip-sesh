package node

import (
	"crypto/ed25519"
	"sync"

	"github.com/georgercarder/same"
)

// read in ~/.ipssh/authorized_keys

// TODO cache authorized keys

var G_SSHMgr = newSSHMgr()

func newSSHMgr() (s *SSHMgr) {
	s = new(SSHMgr)
	return
} // TODO PROPERLY INIT

type SSHMgr struct {
	sync.RWMutex
	// TODO
	pubKeys []*ed25519.PublicKey
}

func (s *SSHMgr) DumpPubKeys() (pks []*ed25519.PublicKey) {
	s.Lock()
	defer s.Unlock()
	for _, pk := range s.pubKeys {
		pks = append(pks, pk)
	}
	return
}

func (s *SSHMgr) IsAuthorized(pk *ed25519.PublicKey) (tf bool) {
	// TODO
	return
}

func getPubKey(domainName string) (pk *ed25519.PublicKey) {
	// TODO
	return
}

func checkPubKeys(
	hash []byte, nonce []byte) (pubKey *ed25519.PublicKey, ok bool) {
	pks := G_SSHMgr.DumpPubKeys()
	// TODO PUT IN THREADS
	for _, pk := range pks {
		if same.Same(hash, Hash(PubKey2Slice(pk), nonce)) {
			pubKey = pk
			ok = true
		}
	}
	return
}
