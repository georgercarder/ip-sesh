package main

import (
	"fmt"

	nd "github.com/georgercarder/ip-sesh/node"
	sg "github.com/georgercarder/ip-sesh/subnet-genie"

	"github.com/ipfs/go-ipfs/core"
)

func main() {
	//go nd.G_Node()

	n := nd.G_Node()
	fmt.Println("debug Identity", n.Identity)
	ps := n.Peerstore.Peers()
	fmt.Println("debug peers", len(ps))
	// fast bootstrap
	sg.FastBootstrap((*core.IpfsNode)(n))
	// announce provide
	go sg.AnnounceProvide((*core.IpfsNode)(n))
	// serve domain
	nd.ServeDomain("test.domain.com")

	// daemon.Listen // TODO
	select {}
}
