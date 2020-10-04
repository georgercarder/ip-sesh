package node

import (
	"crypto/rand"

	gthCrypto "github.com/ethereum/go-ethereum/crypto"
)

func prepareChallenge() (chlg []byte) {
	chlg = genRandToken()
	return
}

func genNonce() (n []byte) {
	n = genRandToken()
	return
}

func genRandToken() (tk []byte) {
	tk = make([]byte, 32)
	rand.Read(tk)
	return
}

func Hash(b ...[]byte) (h []byte) {
	var all []byte
	for _, bb := range b {
		all = append(all, bb...)
	}
	h = gthCrypto.Keccak256(all)
	return
}
