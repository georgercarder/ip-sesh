package main

import (
	"fmt"
	nd "github.com/georgercarder/ipsshd/node"
)

func main() {
	
	g := nd.G_Node()
	fmt.Println("debug", g)
	
	// initializeDomainName // TODO

	// daemon.Listen // TODO
	select{}
}
