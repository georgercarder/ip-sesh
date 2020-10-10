package node

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/georgercarder/lockless-map"
	"github.com/georgercarder/same"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
	fmt.Println("debug SendStop start", domainName)
	hs := m.DomainName2Handshake.Take(domainName)
	if hs != nil {
		fmt.Println("debug SendStop")
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

// TODO

// client

// 1. pub {domName, nonce + H(pubKey, nonce)}

func StartHandshake(domainName string) {
	pubKey := pickPubKey(domainName) // in ssh_mgr
	nonce := genNonce()              // in crypto
	fmt.Println("debug pubKey", pubKey)
	fmt.Println("debug input", Key2Slice(pubKey), nonce)
	hash := Hash(Key2Slice(pubKey), nonce)
	fmt.Println("debug hash", hash)
	hp := &HandshakePacket{DomainName: domainName,
		Nonce: nonce, Hash: hash}
	G_HandshakeMgr.newHandshake(hp, pubKey)
	go publishUntilChallenge(hp)
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
	fmt.Println("debug Slice2HandshakePacket", b)
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
// TODO
func prepareChallengeResponse(
	hp *HandshakePacket) (res *HandshakePacket, err error) {
	// hp does not have domain name (it could but server can cheat),
	// so best to map nonce to domain name internally
	h := G_HandshakeMgr.Nonce2Handshake.Take(string(hp.Nonce))
	if h == nil {
		err = fmt.Errorf("nonce does not map to a handshake.")
		return
	}
	hs := h.(*Handshake)
	G_HandshakeMgr.SendStop(hs.DomainName)
	priv := g_SSHMgr().PubKey2PrivKey.Take(Key2String(hs.PubKey))
	if priv == nil {
		err = fmt.Errorf("PubKey does not map to a priv key.")
		return
	}
	fmt.Println("debug priv", priv)
	// sign
	sig, err := Sign(priv.(*ed25519.PrivateKey), hp.Challenge)
	fmt.Println("debug sig err", sig, err)
	if err != nil {
		// LOG
		return
	}
	// send back to server
	hp.Signature = sig
	res = hp
	return
}

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
	pubKey, ok := checkPubKeys(domainName, hp.Hash, hp.Nonce)
	if !ok {
		// log that pubkey doesn't exist
		return
	}
	fmt.Println("debug pubKey, ok", pubKey, ok)
	s, err := StartStream(context.Background(), pid)
	if err != nil {
		// TODO HANDLE
	}
	sendChallenge(s, hp.Nonce, pubKey)
}

// 4. validate challenge response, put pid in active screens
// TODO
