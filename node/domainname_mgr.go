package node

import (
	"fmt"
	"time"

	"github.com/georgercarder/alerts"
	. "github.com/georgercarder/lockless-map"
	mi "github.com/georgercarder/mod_init"
)

// TODO discovery

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
	// if domName not set, SafelyShutdown emitting error
	// otherwise calls dnsLink w domName
	nn := new(DomNameMgr)
	nn.names = NewLocklessMap() // make(map[string]*DomName)
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
	names LocklessMap // map[string]*DomName
}

type DomName struct {
	name  string
	alert (<-chan interface{})
}

func (d *DomNameMgr) serveDomain(domainName string) {
	dn := &DomName{name: domainName, alert: newDomainAlert(domainName)}
	go handleDomAlerts(dn)
	d.names.Put(domainName, dn)
	Subscribe(domainName, domainName) // topic, alertName
}

func handleDomAlerts(dn *DomName) {
	name := dn.name
	for {
		a := <-dn.alert
		go checkAndRespondToAlert(name, a.([]byte))
	}
}

func newDomainAlert(domainName string) <-chan interface{} {
	sub, err := alerts.G_Alerts().NewSubscription(domainName)
	if err != nil {
		// TODO log
	}
	return sub
}

// TODO SUBSCRIBE TO ALERTS

func (d *DomNameMgr) getDomainNames() (names []string) {
	dumped := d.names.Dump()
	for _, n := range dumped.Keys {
		names = append(names, n.(string))
	}
	return
}
