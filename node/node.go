package node

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs/core"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	
	mi "github.com/georgercarder/mod_init"
)

type IpfsNode core.IpfsNode

func G_Node() (n *IpfsNode) {
	nn, err := modInitializerIpfs.Get()
	if err != nil {
		//LogError.Println("G_Node:", err)
		//reason := err
		//SafelyShutdown(reason)
		return
	}
	return nn.(*IpfsNode)
	return
}

const ModInitTimeout = 3 * time.Second // TODO tune

var modInitializerIpfs = mi.NewModInit(newIpfsNode,
	ModInitTimeout, fmt.Errorf("*IpfsNode init error."))

func newIpfsNode() (n interface{}) { // *IpfsNode
	ncfg := &core.BuildCfg{
		Permanent: true,
		// It is temporary way to signify that node is permanent
		Online:                      true,
		DisableEncryptedConnections: false,
		ExtraOpts: map[string]bool{
			"mplex":  true,
			"pubsub": true,
		},

		Routing: libp2p.DHTClientOption,
	}
	ctx := context.Background()
	node, err := core.NewNode(ctx, ncfg)
	if err != nil {
		//LogError.Println("NewIpfsNode:", err)
		return
	}
	nn := (*IpfsNode)(node)
	err = nn.SetPrivateKey()
	if err != nil {
		//LogError.Println("newIpfsNode:", err)
		return
	}
	n = nn
	return
}

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

func FindPeer(ctx context.Context, pid peer.ID) (pAddrInfo peer.AddrInfo,
	err error) {
	return G_Node().DHT.FindPeer(ctx, pid)
}
