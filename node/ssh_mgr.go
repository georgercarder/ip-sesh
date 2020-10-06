package node

import (
	"crypto/ed25519"
	"fmt"
	"sync"

	mi "github.com/georgercarder/mod_init"
	"github.com/georgercarder/same"
)

func G_SSHMgr() (n *SSHMgr) {
	nn, err := modInitializerSSHMgr.Get()
	if err != nil {
		//LogError.Println("G_Node:", err)
		//reason := err
		//SafelyShutdown(reason)
		return
	}
	return nn.(*SSHMgr)
}

var modInitializerSSHMgr = mi.NewModInit(newSSHMgr,
	ModInitTimeout, fmt.Errorf("*SSHMgr init error."))

// read in ~/.ipssh/authorized_keys

// TODO cache authorized keys

func newSSHMgr() (s interface{}) { //*SSHMgr
	ss := new(SSHMgr)
	s = ss
	return
} // TODO PROPERLY INIT

type SSHMgr struct {
	sync.RWMutex
	// TODO
	privKeys map[string]bool
	pubKeys  map[string]bool
}

func (s *SSHMgr) ImportKeypair(
	priv ed25519.PrivateKey, pub ed25519.PublicKey) (err error) {
	s.Lock()
	defer s.Unlock()
	if priv == nil || pub == nil {
		err = fmt.Errorf("*SSHMgr: keypair must be non-nil.")
		return
	}
	s.privKeys[Key2String(priv)] = true
	s.pubKeys[Key2String(pub)] = true
	return
}

func (s *SSHMgr) DumpPubKeys() (pks []*ed25519.PublicKey) {
	s.Lock()
	defer s.Unlock()
	for s, _ := range s.pubKeys {
		pk := String2PubKey(s)
		pks = append(pks, &pk)
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
	pks := G_SSHMgr().DumpPubKeys()
	// TODO PUT IN THREADS
	for _, pk := range pks {
		if same.Same(hash, Hash(Key2Slice(pk), nonce)) {
			pubKey = pk
			ok = true
		}
	}
	return
}
