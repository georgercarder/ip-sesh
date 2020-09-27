package node

import (
	"context"
	"fmt"
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
			//"pubsub": true,
		},

		Routing: libp2p.DHTClientOption,
	}
	ctx := context.Background()
	node, err := core.NewNode(ctx, ncfg)
	if err != nil {
		//LogError.Println("NewIpfsNode:", err)
		return
	}
	n = (*IpfsNode)(node)
	return
}

func FindPeer(ctx context.Context, pid peer.ID) (pAddrInfo peer.AddrInfo,
	err error) {
	return G_Node().DHT.FindPeer(ctx, pid)
}
