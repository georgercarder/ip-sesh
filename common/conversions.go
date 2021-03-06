package node

import (
	"crypto/ed25519"
	"crypto/sha256"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

func String2CID(key string) (c cid.Cid, err error) {
	h := sha256.New()
	h.Write([]byte(key))
	hashed := h.Sum(nil)
	multihash, err := mh.Encode([]byte(hashed), mh.SHA2_256)
	if err != nil {
		return
	}
	c = cid.NewCidV0(multihash)
	return
}

func Key2Slice(k interface{}) (b []byte) {
	switch k.(type) {
	case *ed25519.PublicKey:
		kk := k.(*ed25519.PublicKey)
		if kk == nil {
			// TODO LOG
			return
		}
		b = []byte(*kk)
		return
	case ed25519.PublicKey:
		kk := k.(ed25519.PublicKey)
		b = []byte(kk)
		return
	case *ed25519.PrivateKey:
		kk := k.(*ed25519.PrivateKey)
		if kk == nil {
			// TODO LOG
			return
		}
		b = []byte(*kk)
		return
	case ed25519.PrivateKey:
		kk := k.(ed25519.PrivateKey)
		b = []byte(kk)
		return
	}
	return
}

func Key2String(k interface{}) (s string) {
	b := Key2Slice(k)
	s = string(b)
	return
}

func String2PubKey(s string) (pk ed25519.PublicKey) {
	b := []byte(s)
	pk = (ed25519.PublicKey)(b)
	return
}

func Slice2PubKey(s []byte) (pk ed25519.PublicKey) {
	pk = (ed25519.PublicKey)(s)
	return
}

func Slice2PrivKey(s []byte) (pk ed25519.PrivateKey) {
	pk = (ed25519.PrivateKey)(s)
	return
}

func String2PrivKey(s string) (pk *ed25519.PrivateKey) {
	p := Slice2PrivKey([]byte(s))
	pk = &p
	return
}

func PubFromPriv(priv ed25519.PrivateKey) (pub ed25519.PublicKey) {
	pub = make([]byte, ed25519.PublicKeySize)
	copy(pub[:], priv[32:]) // see crypto/ed25519 line 91
	return
}

// trims byte(0), byte(0xa) byte(0x20)
func Trim(s string) (ret string) {
	// TODO PUT IN COMMON BUT IS HERE
	// SINCE IS USED TO COMPARE STARTSHELLSESSION
	lastByte := byte(s[len(s)-1])
	if lastByte == byte(0x0) || lastByte == byte(0xa) || lastByte == byte(0x20) {
		return Trim(s[:len(s)-1])
	}
	ret = s
	return

}
