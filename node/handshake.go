package node

import (
	"encoding/json"
	"time"
)

// TODO

// client

// 1. pub {domName, nonce + H(pubKey, nonce)}

func StartHandshake(domainName string) {
	pubKey := getPubKey(domainName) // in ssh_mgr
	nonce := genNonce()             // in crypto
	hash := Hash(pubKey, nonce)
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
}

func (hp *HandshakePacket) Bytes() (b []byte) {
	b, err := json.Marshal(hp)
	if err != nil {
		// TODO LOG ERR
	}
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

// 4. validate challenge response, put pid in active screens
