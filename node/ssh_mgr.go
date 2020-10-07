package node

import (
	"crypto/ed25519"
	"fmt"

	"strconv"
	"strings"
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

func newSSHMgr() (s interface{}) { //*SSHMgr
	ss := new(SSHMgr)
	ss.Domains = newDomains()
	ss.Authorized = newAuthorized()
	ss.setCachedKeys()
	s = ss
	return
} // TODO PROPERLY INIT

func newDomains() (d *domains) {
	d = new(domains)
	d.DomainName2PrivKeys = make(map[string]*keys)
	return
}

func newAuthorized() (a *authorized) {
	a = new(authorized)
	a.DomainName2PubKeys = make(map[string]*keys)
	return
}

type SSHMgr struct {
	sync.RWMutex
	Domains    *domains
	Authorized *authorized
}

type domains struct {
	DomainName2PrivKeys map[string]*keys
}

type authorized struct {
	DomainName2PubKeys map[string]*keys
}

func newKeys() (k *keys) {
	k = new(keys)
	k.M = make(map[string]bool)
	return
}

type keys struct {
	M map[string]bool
}

func (s *SSHMgr) setCachedKeys() {
	sesh_path, err := SESH_Path()
	if err != nil {
		// TODO LOG
		return
	}
	cfg := G_ConfigMgr.Cfg
	// client
	for _, d := range cfg.Client.Domains {
		data, err := SafeFileRead(FSJoin(sesh_path, d.Keyfile))
		if err != nil {
			// TODO LOG
			return
		}
		s.ImportPrivKey(d.Domain, Slice2PrivKey(data))
	}
	// server
	for _, a := range cfg.Server.Authorized {
		data, err := SafeFileRead(FSJoin(sesh_path, a.Keyfile))
		if err != nil {
			// TODO LOG
			return
		}
		s.ImportPubKey(a.Domain, Slice2PubKey(data))
	}
	return
}

type SESH_FILETYPE byte

const (
	UNKNOWN SESH_FILETYPE = iota
	PRIV
	PUB
	CONFIG
)

func identifyFile(fname string) (ft SESH_FILETYPE) {
	split := strings.Split(fname, ".")
	_, err := strconv.Atoi(split[len(split)-1])
	if err == nil { // means last is idx
		split = split[:len(split)-1] // cut off idx
	}
	if len(split) == 1 {
		ft = PRIV
		return
	}
	switch split[len(split)-1] {
	case "pub":
		ft = PUB
		return
	case "config":
		ft = CONFIG
		return

	}
	return // UNKNOWN
}

func (s *SSHMgr) ImportPubKey(
	domainName string, pub ed25519.PublicKey) (err error) {
	s.Lock()
	defer s.Unlock()
	if pub == nil {
		err = fmt.Errorf("*SSHMgr: key must be non-nil.")
		return
	}
	if s.Authorized.DomainName2PubKeys[domainName] == nil {
		s.Authorized.DomainName2PubKeys[domainName] = newKeys()
	}
	s.Authorized.DomainName2PubKeys[domainName].M[Key2String(pub)] = true
	return
}

func (s *SSHMgr) ImportPrivKey(
	domainName string, priv ed25519.PrivateKey) (err error) {
	s.Lock()
	defer s.Unlock()
	if priv == nil {
		err = fmt.Errorf("*SSHMgr: key must be non-nil.")
		return
	}
	if s.Domains.DomainName2PrivKeys[domainName] == nil {
		s.Domains.DomainName2PrivKeys[domainName] = newKeys()
	}
	s.Domains.DomainName2PrivKeys[domainName].M[Key2String(priv)] = true
	return
}

func (s *SSHMgr) DumpPubKeys(domainName string) (pks []*ed25519.PublicKey) {
	s.Lock()
	defer s.Unlock()
	pubKeys := s.Authorized.DomainName2PubKeys[domainName]
	for s, _ := range pubKeys.M {
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

func checkPubKeys(domainName string,
	hash, nonce []byte) (pubKey *ed25519.PublicKey, ok bool) {
	pks := G_SSHMgr().DumpPubKeys(domainName)
	// TODO PUT IN THREADS
	for _, pk := range pks {
		if same.Same(hash, Hash(Key2Slice(pk), nonce)) {
			pubKey = pk
			ok = true
		}
	}
	return
}
