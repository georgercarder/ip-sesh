package node

import (
	"crypto/ed25519"
)

// TODO need consts.go

const Home_Dir = "~/"

var SESH_Path = FSJoin(Home_Dir, ".sesh")

func GenerateAndSaveKeypair(filename string) (err error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	// nil setting defaults to crypto/rand
	if err != nil {
		// TODO LOG
		return
	}
	privPath := FSJoin(SESH_Path, filename)
	err = SafeFileWrite(privPath, Key2Slice(priv))
	if err != nil {

	}
	pubPath := FSJoin(SESH_Path, filename+".pub")
	err = SafeFileWrite(pubPath, Key2Slice(pub))
	if err != nil {

	}
	err = G_SSHMgr().ImportKeypair(priv, pub)
	if err != nil {
		// TODO LOG
	}
	return
}
