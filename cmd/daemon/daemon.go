package main

import (
	nd "github.com/georgercarder/ipsshd/node"
)

func main() {
	go nd.G_Node()
	
	// daemon.Listen // TODO
	select{}
}
