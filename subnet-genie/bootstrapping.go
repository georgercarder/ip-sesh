package subnet_genie

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	nd "github.com/georgercarder/ip-sesh/node"

	"github.com/ipfs/go-ipfs/core"
	"github.com/libp2p/go-libp2p-core/peer"
)

func FastBootstrap(n *core.IpfsNode) {
	numPeers := 0
	for numPeers < 1000 {
		ps := n.Peerstore.Peers()
		numPeers = len(ps)
		fmt.Println("debug peers", len(ps))
		time.Sleep(100 * time.Millisecond)

		go func() {
			dht := n.DHT
			rval := make([]byte, 32)
			rand.Read(rval)
			ctx, cancel := context.WithTimeout(
				context.Background(), 2*time.Second)
			defer cancel()
			v, err := dht.GetValue(ctx, string(rval))

			if err != nil {
				//fmt.Println("debug err", err)
			}
			fmt.Println("debug v", v)
		}()
	}
}

func AnnounceProvide(n *core.IpfsNode) {
	for { // an interval
		key, err := nd.String2CID("/ip-sesh/0.0.1")
		if err != nil {
			fmt.Println("debug conv err", err)
		}
		ctx, cancel := context.WithTimeout(
			context.Background(), 10*time.Second)
		defer cancel()
		err = n.Routing.Provide(ctx, key, true)
		if err != nil {
			fmt.Println("debug provide err", err)
		}
		fmt.Println("debug provided")
		time.Sleep(5 * time.Minute)
	}
}

// the point is to connect with providers of "/ip-sesh/0.0.1"
// to propagate the pubsub pubs to these "providers"
func JoinProviders(n *core.IpfsNode) {
	key, err := nd.String2CID("/ip-sesh/0.0.1")
	if err != nil {
		fmt.Println("debug conv err", err)
	}
	numProvs := 1
	foundOne := false
	for !foundOne && numProvs < 1024 { // TODO make max const
		ctx, cancel := context.WithTimeout(
			context.Background(), 10*time.Second)
		defer cancel()
		pchan := n.DHT.FindProvidersAsync(ctx, key, numProvs)
		numProvs *= 2
		ct := 0
		for ct < numProvs { // TODO PUT IN TIMEOUT
			ct++
			p := <-pchan
			//fmt.Println("debug Provider found", p)
			if len(p.Addrs) < 1 {
				continue
			}
			go func(pp peer.AddrInfo) {
				ctx2, cancel2 := context.WithTimeout(
					context.Background(), 10*time.Second)
				defer cancel2()
				pr, err := nd.FindPeer(ctx2, pp.ID)
				if err != nil {
					fmt.Println("debug FindPeer", err)
				}
				fmt.Println("debug FindPeer done", pr.ID)
				foundOne = true // this is janky FIXME
			}(p)
		}
	}
}
