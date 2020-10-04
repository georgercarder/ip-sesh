package node

import (
	"fmt"
	"sync"
	"time"

	"github.com/georgercarder/alerts"
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
	nn.names = make(map[string]*DomName)
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
	names map[string]*DomName
}

type DomName struct {
	name  string
	alert (<-chan interface{})
}

func (d *DomNameMgr) serveDomain(domainName string) {
	d.Lock()
	defer d.Unlock()
	dn := &DomName{name: domainName, alert: newDomainAlert(domainName)}
	go handleDomAlerts(dn)
	d.names[domainName] = dn
	Subscribe(domainName, domainName) // topic, alertName
}

func handleDomAlerts(dn *DomName) {
	//name := dn.name
	for {
		a := <-dn.alert
		go checkAndRespondToAlert(a.([]byte))
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
	d.RLock()
	defer d.Unlock()
	for n, _ := range d.names {
		names = append(names, n)
	}
	return
}
