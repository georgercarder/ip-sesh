package node

import (
	"fmt"
	"time"
	
	mi "github.com/georgercarder/mod_init"
)

func G_DomNameMgr() (m *DomNameMgr) {
	mm, err := modInitializerDomNameMgr.Get()
	if err != nil {
		//LogError.Println("G_DomNameMgr:", err)
		//reason := err
		//SafelyShutdown(reason)
		return
	}
	return mm.(*DomNameMgr)
}

const DomNameMgrInitTimeout = 3 * time.Second // TODO tune

var modInitializerDomNameMgr = mi.NewModInit(newDomNameMgr,
	DomNameMgrInitTimeout, fmt.Errorf("*DomNameMgr init error."))

type DomNameMgr struct {
	// TODO
}

func newDomNameMgr() (n interface{}) { // *DomNameMgr
	fmt.Println("debug newDomeNameMgr")
	// TODO
	// checks /etc/ipsshd/ipsshd.conf
	// if domName not set, SafelyShutdown emitting error
	// otherwise calls dnsLink w domName

	return
}
