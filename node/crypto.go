package node

import (
	"crypto/rand"
)

func prepareChallenge() (chlg []byte, err error) {
	randToken, err := genRandToken()
	if err != nil {
		return
	}
	chlg = randToken
	return
}

func genRandToken() (tk []byte, err error) {
	tk = make([]byte, 32)
	_, err = rand.Read(tk)
	return
}
