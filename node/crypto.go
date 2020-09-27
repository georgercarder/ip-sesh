package node

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"sync"
)

var g_NodeData = new(NodeData)

type NodeData struct {
	sync.RWMutex
	privateKey *rsa.PrivateKey
}

func (n *IpfsNode) SetPrivateKey() (err error) {
	raw, err := n.PrivateKey.Raw()
	if err != nil {
		return
	}
	pk, err := x509.ParsePKCS1PrivateKey(raw)
	if err != nil {
		return
	}
	// note: g_NodeData is locked by caller function
	g_NodeData.privateKey = pk
	return
}

func (n *IpfsNode) PublicKey() (p *rsa.PublicKey, err error) {
	g_NodeData.RLock()
	defer g_NodeData.RUnlock()
	p = g_NodeData.privateKey.Public().(*rsa.PublicKey)
	return
}

func (n *IpfsNode) DecryptOAEPRSA(bMessage []byte) (p []byte, err error) {
	g_NodeData.RLock()
	defer g_NodeData.RUnlock()
	return DecryptOAEPRSA(bMessage, g_NodeData.privateKey)
}

func DecryptOAEPRSA(bMessage []byte, pk *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	r := rand.Reader
	return rsa.DecryptOAEP(h, r, pk, bMessage, nil)
}

func EncryptOAEPRSA(bMessage []byte, pk *rsa.PublicKey) ([]byte, error) {
	h := sha256.New()
	r := rand.Reader
	return rsa.EncryptOAEP(h, r, pk, bMessage, nil)
}

func prepareChallenge(pk *rsa.PublicKey) (chlg []byte, err error) {
	randToken, err := genRandToken()
	if err != nil {
		return
	}
	encryptedToken, err := EncryptOAEPRSA(randToken, pk)
	chlg = encryptedToken // for explicit readability
	return
}

func genRandToken() (tk []byte, err error) {
	// TODO
	return
}
