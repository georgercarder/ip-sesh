package node

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/georgercarder/alerts"
	mi "github.com/georgercarder/mod_init"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// TODO put in own module.
// this is generic and isolated to be own module.

func G_pubSub() (p *pubSub) {
	pp, err := modInitializerPubSub.Get()
	if err != nil {
		//LogError.Println("G_pubSubCH:", err)
		//reason := err
		//SafelyShutdown(reason)
		return
	}
	return pp.(*pubSub)
}

var modInitializerPubSub = mi.NewModInit(newPubSub,
	ModInitTimeout, fmt.Errorf("*pubSub init error."))

func PubSubTopicHashed(str string) (r string) {
	keyAsCID, err := String2CID(str)
	if err != nil {
		//LogError.Println("newPubSub:", err)
	}
	r = keyAsCID.String()
	return
}

func newPubSub() (p interface{}) { // *pubSub
	pp := new(pubSub)
	pp.Lock()
	defer pp.Unlock()
	pp.M = make(map[string]*subNAlert)
	p = pp
	return
}

func handleSubscription(sNa *subNAlert) {
	fmt.Println("debug handleSubscription", sNa)
	for {
		msg, err := sNa.s.Next(context.Background())
		fmt.Println("debug readSubs:", msg, err)
		b, err := json.Marshal(msg)
		if err != nil {
			//LogError.Println("handleSubscription:", err)
			continue
		}
		go alerts.G_Alerts().SendAlert(sNa.alertName, b)
	}
}

type pubSub struct {
	sync.RWMutex
	M map[string]*subNAlert
}

type subNAlert struct {
	s         *pubsub.Subscription
	alertName string
}

func PubSubGetTopics() (ls []string) {
	n := G_Node()
	return n.PubSub.GetTopics()
}

func Publish(topic string, data []byte) (err error) {
	return publish(topic, data)
}

func publish(topic string, data []byte) (err error) {
	key := PubSubTopicHashed(topic)
	fmt.Println("debug Publish", topic, key, data)
	n := G_Node()
	return n.PubSub.Publish(key, data)
}

func Subscribe(topic, alertName string) (err error) {
	key := PubSubTopicHashed(topic)
	sub, err := subscribe(key)
	if err != nil {
		//LogError.Println("newPubSub:", err)
		return
	}
	sNa := &subNAlert{s: sub,
		alertName: alertName}
	G_pubSub().M[topic] = sNa
	go handleSubscription(sNa)
	return
}

func subscribe(topic string) (*pubsub.Subscription, error) {
	fmt.Println("debug subscribe", topic)
	n := G_Node()
	return n.PubSub.Subscribe(topic)
}
