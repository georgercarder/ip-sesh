package main

import (
	"bufio"

	"fmt"
	"os"

	nd "github.com/georgercarder/ip-sesh/node"
	sg "github.com/georgercarder/ip-sesh/subnet-genie"

	"github.com/ipfs/go-ipfs/core"
)

func main() {
	/*err := nd.GenerateAndSaveKeypair("cats")
	if err != nil {
		fmt.Println("debug key err", err)
	}*/
	fmt.Println("client")
	fmt.Println("initializing node ...")
	n := nd.G_Node()
	// fast bootstrap
	sg.FastBootstrap((*core.IpfsNode)(n))
	ps := n.Peerstore.Peers()
	fmt.Println("peers", len(ps))
	go sg.JoinProviders((*core.IpfsNode)(n))
	fmt.Println("Press ENTER for demo.")
	getCharReader := bufio.NewReader(os.Stdin)
	_, err := getCharReader.ReadString('\n')
	if err != nil {
		fmt.Println("debug error", err)
	}
	nd.StartHandshake("test.domain.com")

	select {}
}
