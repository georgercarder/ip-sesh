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
	go func(){
		time.Sleep(1*time.Second)
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
		nd.FindPeer(ctx2, p.ID)
		fmt.Println("debug FindPeer done")
	}()

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

		/*numPeers := 0
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
		}*/


			/*
			cid, err := nd.String2CID("TEST")
			if err != nil{
				fmt.Println("debug err", err)
			}
			fmt.Println("debug cid", cid)
			pid, err := peer.FromCid(cid)
			if err != nil{
				fmt.Println("debug err", err)
			}
			p, err := nd.FindPeer(context.Background(), pid)
			if err != nil{
				fmt.Println("debug err", err)
			}
			fmt.Println("peer found", p)
			*/
