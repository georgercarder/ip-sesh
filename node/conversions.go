package node

// CONVERSIONS TODO PUT IN COMMON

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
	case *ed25519.PrivateKey:
		kk := k.(*ed25519.PrivateKey)
		if kk == nil {
			// TODO LOG
			return
		}
		b = []byte(*kk)
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
