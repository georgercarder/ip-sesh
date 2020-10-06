package node

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/georgercarder/same"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

var G_HandshakeMgr = newHandshakeMgr() // TODO USE MODINIT

func newHandshakeMgr() (h *HandshakeMgr) {
	h = new(HandshakeMgr)
	h.DomainName2Handshake = make(map[string]*Handshake)
	return
}

type HandshakeMgr struct {
	sync.RWMutex
	DomainName2Handshake map[string]*Handshake
}

func (m *HandshakeMgr) newHandshake(domainName string) {
	m.Lock()
	defer m.Unlock()
	hs := new(Handshake)
	hs.StopChnl = make(chan bool)
	m.DomainName2Handshake[domainName] = hs
}

func (m *HandshakeMgr) SendStop(domainName string) {
	fmt.Println("debug SendStop", domainName)
	if m.DomainName2Handshake[domainName] != nil {
		fmt.Println("debug SendStop")
		stop := true
		m.DomainName2Handshake[domainName].StopChnl <- stop
	}
}

func (m *HandshakeMgr) Stop(domainName string) <-chan bool {
	m.Lock()
	hs := m.DomainName2Handshake[domainName]
	m.Unlock()
	return hs.StopChnl
}

type Handshake struct {
	DomainName string
	StopChnl   (chan bool)
}

// TODO

// client

// 1. pub {domName, nonce + H(pubKey, nonce)}

func StartHandshake(domainName string) {
	pubKey := getPubKey(domainName) // in ssh_mgr
	nonce := genNonce()             // in crypto
	hash := Hash(PubKey2Slice(pubKey), nonce)
	hp := &HandshakePacket{DomainName: domainName,
		Nonce: nonce, Hash: hash}
	G_HandshakeMgr.newHandshake(domainName)
	go publishUntilChallenge(hp)
}

const publishInterval = 1 * time.Second // TODO put elsewhere and tune

type HandshakePacket struct {
	DomainName string
	Nonce      []byte
	Hash       []byte
	Challenge  []byte
	PubKey     *ed25519.PublicKey
}

func (hp *HandshakePacket) Bytes() (b []byte) {
	b, err := json.Marshal(hp)
	if err != nil {
		// TODO LOG ERR
	}
	return
}

func Slice2HandshakePacket(b []byte) (hp *HandshakePacket, err error) {
	hp = new(HandshakePacket)
	err = json.Unmarshal(b, &hp)
	return
}

func publishUntilChallenge(hp *HandshakePacket) {
	pubCH := make(chan bool)
	go func() {
		for {
			pubCH <- true
			time.Sleep(publishInterval)
		}
	}()
	stopper := G_HandshakeMgr.Stop(hp.DomainName)
	fmt.Println("debug stopper", stopper)
	for {
		select {
		case <-pubCH:
			Publish(hp.DomainName, hp.Bytes())
		case <-stopper:
			return
		}
	}
}

// 3. respond to challenge
// put this in
//G_HandshakeMgr.SendStop(domainName)
// server

// 2. from alert, scan pubKeys w nonce, compare to H,
//    if has pubKey, send challenge to client
func checkAndRespondToAlert(domainName string, a []byte) {
	m := new(pubsub.Message)
	err := json.Unmarshal(a, &m)
	if err != nil {
		// TODO LOG
	}
	hp := new(HandshakePacket)
	err = json.Unmarshal(m.Data, &hp)
	if err != nil {
		// TODO LOG ERR
	}
	if !same.Same(domainName, hp.DomainName) {
		// LOG
		// means client is trying to cheat
		return
	}
	pid, err := peer.IDFromBytes(m.From)
	if err != nil {
		/// TODO LOG ERR
	}
	fmt.Println("debug pid", pid)
	pubKey, ok := checkPubKeys(hp.Hash, hp.Nonce)
	if !ok {
		// log that pubkey doesn't exist
		return
	}
	fmt.Println("debug pubKey, ok", pubKey, ok)
	s, err := StartStream(context.Background(), pid)
	if err != nil {
		// TODO HANDLE
	}
	go sendChallenge(s, pubKey)
}

// 4. validate challenge response, put pid in active screens
