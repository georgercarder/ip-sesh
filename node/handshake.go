package node

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// TODO

// client

// 1. pub {domName, nonce + H(pubKey, nonce)}

func StartHandshake(domainName string) {
	pubKey := getPubKey(domainName) // in ssh_mgr
	nonce := genNonce()             // in crypto
	hash := Hash(PubKey2Slice(pubKey), nonce)
	hp := &HandshakePacket{DomainName: domainName,
		Nonce: nonce, Hash: hash}
	// TODO register w Handshake_Mgr
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
	for {
		select {
		case <-pubCH:
			Publish(hp.DomainName, hp.Bytes())
			/*case <-G_Handshake_Mgr.Stop(hp.DomainName):
			return*/ // TODO
		}
	}
}

// 3. respond to challenge

// server

// 2. from alert, scan pubKeys w nonce, compare to H,
//    if has pubKey, send challenge to client
func checkAndRespondToAlert(a []byte) {
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
