package node

import (
	"encoding/json"
	"sync"
)

var G_ConfigMgr = newConfigMgr() // TODO use mod_init

func newConfigMgr() (c *configMgr) {
	cc := new(configMgr)
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
		SaveConfig(cc.Cfg) // puts a template for user to fill out
	}
	c = cc
	return
}

type configMgr struct {
	sync.RWMutex
	Cfg *Config
}

type Config struct {
	Todo string `json:"todo"` // TODO
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
