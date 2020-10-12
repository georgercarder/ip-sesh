package node

import (
	"encoding/json"
	"sync"

	. "github.com/georgercarder/lockless-map"
)

var G_ConfigMgr = newConfigMgr() // TODO use mod_init

func newConfigMgr() (c *configMgr) {
	cc := new(configMgr)
	cc.Lock()
	defer cc.Unlock()
	cc.Cfg = new(Config)
	sesh_config_path, err := SESH_Config_Path()
	if err != nil {
		// TODO LOG
		return
	}
	if FileExists(sesh_config_path) {
		data, err := SafeFileRead(sesh_config_path)
		if err != nil {
			// TODO LOG
			return
		}
		err = json.Unmarshal(data, &cc.Cfg)
	} else {
		cc.Cfg.Client = new(Client)
		sampleDomain := &DomainNKeyfile{Domain: "example.com",
			Keyfile: "example.pub"}
		cc.Cfg.Client.Domains = append(cc.Cfg.Client.Domains,
			sampleDomain)
		cc.Cfg.Server = new(Server)
		sampleAuthorized := &DomainNKeyfile{Domain: "example2.com",
			Keyfile: "example2.pub"}
		cc.Cfg.Server.Authorized = append(cc.Cfg.Server.Authorized,
			sampleAuthorized)
		SaveConfig(cc.Cfg) // puts a template for user to fill out
	}
	G_ValidClientDomains = NewLocklessMap()
	for _, d := range cc.Cfg.Client.Domains {
		G_ValidClientDomains.Put(d.Domain, true)
	}
	c = cc
	return
}

type configMgr struct {
	sync.RWMutex
	Cfg *Config
}

type Config struct {
	Client *Client `json:"client"`
	Server *Server `json:"server"`
	// TODO eventually specify user. But for now session
	// has root privileges.
}

type Client struct {
	Domains []*DomainNKeyfile `json:"domains"`
}

var G_ValidClientDomains LocklessMap

func ClientDomainIsValid(domain string) (ok bool) {
	d := G_ValidClientDomains.Take(domain)
	if d == nil {
		return
	}
	ok = d.(bool)
	return
}

type Server struct {
	Authorized []*DomainNKeyfile `json:"authorized"`
}

type DomainNKeyfile struct {
	Domain  string `json:"domain"`
	Keyfile string `json:"keyfile"`
}

func SaveConfig(conf interface{}) (err error) {
	b, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		// TODO LOG ERR
		return
	}
	sesh_config_path, err := SESH_Config_Path()
	if err != nil {
		// TODO LOG
		return
	}
	return SafeFileWrite(sesh_config_path, b)
}

// TODO getters
