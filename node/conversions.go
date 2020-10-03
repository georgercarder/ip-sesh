package node

// CONVERSIONS TODO PUT IN COMMON

import (
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
