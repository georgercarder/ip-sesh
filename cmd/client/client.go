package main

import (
	"bufio"
	"context"
	"crypto/rand"

	"fmt"
	"os"
	"time"

	nd "github.com/georgercarder/ip-sesh/node"
	sg "github.com/georgercarder/ip-sesh/subnet-genie"

	"github.com/ipfs/go-ipfs/core"
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
	go sg.JoinProviders((*core.IpfsNode)(n))
	for {
		fmt.Println("Press ENTER for demo.")
		getCharReader := bufio.NewReader(os.Stdin)
		_, err := getCharReader.ReadString('\n')
		if err != nil {
			fmt.Println("debug error", err)
		}
		nd.StartHandshake("test.domain.com")
	}

	select {}
}
