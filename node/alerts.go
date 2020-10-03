package node

import (
	"fmt"
	"sync"

	mi "github.com/georgercarder/mod_init"
)

func G_Alerts() (a *Alerts) {
	aa, err := modInitializeAlerts.Get()
	if err != nil {
		//LogError.Println("G_Alerts:", err)
		//reason := err
		//SafelyShutdown(reason)
		return
	}
	return aa.(*Alerts)
}

var modInitializeAlerts = mi.NewModInit(NewAlertsHub,
	ModInitTimeout, fmt.Errorf("*Alerts init error"))

type Alerts struct {
	sync.RWMutex
	Name2Chan map[string]*ChanString
}

func NewAlertsHub() (a interface{}) { // *Alerts
	aa := new(Alerts)
	aa.Init()
	a = aa
	return
}

func (a *Alerts) SendAlert(chanName string, data string) {
	if a.Name2Chan[chanName] == nil {
		a.Lock()
		channel := make(chan string)
		a.newChanLocked(chanName, channel)
		a.Unlock()
	}
	a.RLock()
	a.Name2Chan[chanName].CH <- data
	a.RUnlock()
	return
}

func (a *Alerts) Init() {
	a.Lock()
	defer a.Unlock()
	a.Name2Chan = make(map[string]*ChanString)
}

func (a *Alerts) newChanLocked(chanName string, CH chan string) (er error) {
	// put type deduction for different chans! // TODO
	// for now we only need ChanString
	ch := new(ChanString)
	ch.Init(CH)
	a.Name2Chan[chanName] = ch
	return
}

func (a *Alerts) NewSubscription(chanName string) (subCH <-chan string,
	er error) {
	a.Lock()
	defer a.Unlock()
	if a.Name2Chan[chanName] == nil {
		channel := make(chan string)
		a.newChanLocked(chanName, channel)
	}
	// put type deduction for different chans! // TODO
	// for now we only need ChanString
	ch := a.Name2Chan[chanName]
	subCH = ch.NewSubscription()
	return
}

type Chan interface {
	NewSubscription()
}

//
type ChanString struct {
	sync.RWMutex
	CH            (chan string)
	Subscriptions [](chan string)
}

func (c *ChanString) Init(CH chan string) {
	c.Lock()
	defer c.Unlock()
	c.CH = CH
	c.fanout()
	return
}

// this is effectively creates a fanout on ChanString.CH
func (c *ChanString) NewSubscription() <-chan string {
	c.Lock()
	subscription := make(chan string)
	c.Subscriptions = append(c.Subscriptions, subscription)
	c.Unlock()
	return subscription
}

func (c *ChanString) fanout() {
	go func() {
		for {
			select {
			case s := <-c.CH:
				c.RLock()
				for i := 0; i < len(c.Subscriptions); i++ {
					c.Subscriptions[i] <- s
				}
				c.RUnlock()
			}
		}
	}()
	return
}

//
const (
	PubSub_A = "PubSub"
)
