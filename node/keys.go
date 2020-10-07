package node

import (
	"crypto/ed25519"
	//"fmt"
	"os/user"
)

func Home_Dir() (hd string, err error) {
	usr, err := user.Current()
	if err != nil {
		// TODO LOG
		return
	}
	hd = usr.HomeDir
	return
} // specific to Linux. does not support Windows

func SESH_Path() (sp string, err error) {
	hd, err := Home_Dir()
	if err != nil {
		// TODO LOG
		return
	}
	sp = FSJoin(hd, ".sesh")
	return
}

func SESH_Config_Path() (cp string, err error) {
	sp, err := SESH_Path()
	if err != nil {
		// TODO LOG
		return
	}
	cp = FSJoin(sp, "sesh.config")
	return
}

func GenerateAndSaveKeypair(filename string) (err error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	// nil setting defaults to crypto/rand
	if err != nil {
		// TODO LOG
		return
	}
	sesh_path, err := SESH_Path()
	if err != nil {
		// TODO LOG
		return
	}
	privPath := FSJoin(sesh_path, filename)
	err = SafeFileWrite(privPath, Key2Slice(priv))
	if err != nil {
		// TODO LOG
		return
	}
	pubPath := FSJoin(sesh_path, filename+".pub")
	err = SafeFileWrite(pubPath, Key2Slice(pub))
	if err != nil {
		// TODO LOG
		return
	}
	err = G_SSHMgr().ImportKeypair(priv, pub)
	if err != nil {
		// TODO LOG
	}
	return
}
