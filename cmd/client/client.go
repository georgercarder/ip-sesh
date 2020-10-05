package main 

import (
	"bufio"
	"context"
	"crypto/rand"

	"fmt"
	"os"
	"time"

	nd "github.com/georgercarder/ip-sesh/node"

	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	fmt.Println("debug client")
	n := nd.G_Node()
	numPeers := 0
	// fast bootstrap
	for numPeers < 1000 {
		ps := n.Peerstore.Peers()
		numPeers = len(ps)
		fmt.Println("debug peers", len(ps))
		time.Sleep(100*time.Millisecond)

		go func() {
			dht := n.DHT
			rval := make([]byte, 32)
			rand.Read(rval)
			ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
			defer cancel()
			v, err := dht.GetValue(ctx, string(rval))

			if err != nil{
				//fmt.Println("debug err", err)
			}
			fmt.Println("debug v", v)
		}()
	}
	go joinProviders()
	for {
		fmt.Println("Press ENTER for demo.")
		getCharReader := bufio.NewReader(os.Stdin)
		_, err := getCharReader.ReadString('\n')
		if err != nil {
			fmt.Println("debug error", err)
		}
		nd.StartHandshake("test.domain.com")
	}

	select{}
}

// the point is to connect with providers of "/ip-sesh/0.0.1"
// to propagate the pubsub pubs to these "providers"
func joinProviders() {
	n := nd.G_Node()
	key, err := nd.String2CID("/ip-sesh/0.0.1")
	if err != nil {
		fmt.Println("debug conv err", err)
	}
	numProvs := 1
	foundOne := false
	for !foundOne {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		pchan := n.DHT.FindProvidersAsync(ctx, key, numProvs)
		numProvs *= 2 
		ct := 0
		for ct < 100 {// TODO PUT IN TIMEOUT
			ct++
			p := <-pchan
			//fmt.Println("debug Provider found", p)
			if len(p.Addrs) < 1 {
				continue
			}
			go func(pp peer.AddrInfo) {
				ctx2, cancel2 := context.WithTimeout(context.Background(), 10 * time.Second)
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
