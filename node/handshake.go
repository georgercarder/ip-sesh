package node

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"

	"github.com/georgercarder/same"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// client

// 1. pub {domName, nonce + H(pubKey, nonce)}

func StartHandshake(domainName string) {
	pubKey := pickPubKey(domainName) // in ssh_mgr
	nonce := genNonce()              // in crypto
	//fmt.Println("debug pubKey", pubKey)
	//fmt.Println("debug input", Key2Slice(pubKey), nonce)
	hash := Hash(Key2Slice(pubKey), nonce)
	//fmt.Println("debug hash", hash)
	hp := &HandshakePacket{DomainName: domainName,
		Nonce: nonce, Hash: hash}
	G_HandshakeMgr.newHandshake(hp, pubKey)
	go publishUntilChallenge(hp)
}

// 3. respond to challenge
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
	//fmt.Println("debug priv", priv)
	// sign
	sig, err := Sign(priv.(*ed25519.PrivateKey), hp.Challenge)
	//fmt.Println("debug sig err", sig, err)
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
		return
	}
	hp := new(HandshakePacket)
	err = json.Unmarshal(m.Data, &hp)
	if err != nil {
		// TODO LOG ERR
		return
	}
	if !same.Same(domainName, hp.DomainName) {
		// LOG
		// means client is trying to cheat
		return
	}
	pid, err := peer.IDFromBytes(m.From)
	if err != nil {
		/// TODO LOG ERR
		return
	}
	//fmt.Println("debug pid", pid)
	pubKey, ok := checkPubKeys(domainName, hp.Hash, hp.Nonce)
	if !ok {
		// log that pubkey doesn't exist
		return
	}
	//fmt.Println("debug pubKey, ok", pubKey, ok)
	s, err := StartStream(context.Background(), pid)
	if err != nil {
		// TODO HANDLE
		return
	}
	sendChallenge(s, hp.Nonce, pubKey)
	defer s.Reset()
}

// 4. validate challenge response, put pid in active screens
// see // 3 of stream_handler.go
// func checkHandshakeResponse(s network.Stream) (err error)
