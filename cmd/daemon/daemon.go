package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
	nd "github.com/georgercarder/ip-sesh/node"
)

func main() {
	//go nd.G_Node()

	n := nd.G_Node()
	fmt.Println("debug Identity", n.Identity)
	ps := n.Peerstore.Peers()
	fmt.Println("debug peers", len(ps))
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
	// announce provide
	go announceProvide()
	// serve domain
	nd.ServeDomain("test.domain.com")

	// daemon.Listen // TODO
	select {}
}

func announceProvide() {
	n := nd.G_Node()
	for { // an interval
		key, err := nd.String2CID("/ip-sesh/0.0.1")
		if err != nil {
			fmt.Println("debug conv err", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		err = n.Routing.Provide(ctx, key, true)
		if err != nil {
			fmt.Println("debug provide err", err)
		}
		fmt.Println("debug provided")
		time.Sleep(5*time.Minute)
	}
}
