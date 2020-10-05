package main 

import (
	"bufio"
	"context"
//	"crypto/rand"
	"fmt"
	"os"
	"time"

	nd "github.com/georgercarder/ip-sesh/node"

//	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	fmt.Println("debug client")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	pchan := n.DHT.FindProvidersAsync(ctx, key, 1)
	p := <-pchan
	fmt.Println("debug Provider found", p)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 4 * time.Second)
	defer cancel2()
	pp, err := nd.FindPeer(ctx2, p.ID)
	if err != nil {
		fmt.Println("debug FindPeer", err)
	}
	fmt.Println("debug FindPeer done", pp)
}
