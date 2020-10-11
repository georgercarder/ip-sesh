package node

import (
	"crypto/ed25519"
	"encoding/json"
	"time"

	. "github.com/georgercarder/lockless-map"
)

var G_HandshakeMgr = newHandshakeMgr() // TODO USE MODINIT

func newHandshakeMgr() (h *HandshakeMgr) {
	h = new(HandshakeMgr)
	h.DomainName2Handshake = NewLocklessMap()
	h.Nonce2Handshake = NewLocklessMap()
	return
}

type HandshakeMgr struct {
	DomainName2Handshake LocklessMap // map[string]*Handshake
	Nonce2Handshake      LocklessMap // map[string]*Handshake
}

func (m *HandshakeMgr) newHandshake(
	hp *HandshakePacket, pubKey *ed25519.PublicKey) {
	hs := &Handshake{
		DomainName: hp.DomainName,
		Nonce:      hp.Nonce,
		Challenge:  hp.Challenge,
		PubKey:     pubKey,
		LastTouch:  time.Now()}
	hs.StopChnl = make(chan bool)
	m.DomainName2Handshake.Put(hp.DomainName, hs)
	m.Nonce2Handshake.Put(string(hp.Nonce), hs)
}

func (m *HandshakeMgr) SendStop(domainName string) {
	hs := m.DomainName2Handshake.Take(domainName)
	if hs != nil {
		stop := true
		hhs := hs.(*Handshake)
		hhs.StopChnl <- stop
	}
}

func (m *HandshakeMgr) Stop(domainName string) <-chan bool {
	hs := m.DomainName2Handshake.Take(domainName)
	if hs == nil {
		// TODO LOG
		// this should not!! happen
		return make(chan bool) // empty chnl
	}
	hhs := hs.(*Handshake)
	return hhs.StopChnl
}

type Handshake struct {
	DomainName string
	Nonce      []byte
	Challenge  []byte
	PubKey     *ed25519.PublicKey
	StopChnl   (chan bool)
	LastTouch  time.Time
}

const publishInterval = 1 * time.Second // TODO put elsewhere and tune

type HandshakePacket struct {
	DomainName string
	Nonce      []byte
	Hash       []byte
	Challenge  []byte
	Signature  []byte
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
	for {
		select {
		case <-pubCH:
			Publish(hp.DomainName, hp.Bytes())
		case <-stopper:
			return
		}
	}
}
