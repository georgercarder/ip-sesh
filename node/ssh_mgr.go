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
	pks := G_SSHMgr().DumpPubKeys()
	// TODO PUT IN THREADS
	for _, pk := range pks {
		if same.Same(hash, Hash(PubKey2Slice(pk), nonce)) {
			pubKey = pk
			ok = true
		}
	}
	return
}
