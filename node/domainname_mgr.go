package node

import (
	"fmt"
	"sync"
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

func newDomNameMgr() (n interface{}) { // *DomNameMgr
	fmt.Println("debug newDomNameMgr")
	// TODO
	// checks ~/.ipsshd/ipsshd.conf
	// if domName not set, SafelyShutdown emitting error
	// otherwise calls dnsLink w domName
	nn := new(DomNameMgr)
	nn.names = make(map[string]bool)
	n = nn
	return
}

func ServeDomain(domainName string) {
	G_DomNameMgr().serveDomain(domainName)
}

func GetDomainNames() (names []string) {
	return G_DomNameMgr().getDomainNames()
}

type DomNameMgr struct {
	sync.RWMutex
	names map[string]bool
}

func (d *DomNameMgr) serveDomain(domainName string) {
	d.Lock()
	defer d.Unlock()
	d.names[domainName] = true
	Subscribe(domainName)
}

func (d *DomNameMgr) getDomainNames() (names []string) {
	d.RLock()
	defer d.Unlock()
	for n, _ := range d.names {
		names = append(names, n)
	}
	return
}
