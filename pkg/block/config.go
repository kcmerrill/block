package block

import (
	"github.com/ghodss/yaml"
	"github.com/kcmerrill/common.go/config"
)

// Config ...
type Config struct {
	Overrides map[string]string `yaml:"overrides"`
	Boost     map[string]int    `yaml:"boost"`
}

func (b *Block) config() {
	c := &Config{
		Boost:     make(map[string]int),
		Overrides: make(map[string]string),
	}
	configContents, configErr := config.Home("block", ".yml")
	if configErr == nil {
		unmarshalErr := yaml.Unmarshal([]byte(configContents), &c)
		if unmarshalErr == nil {
			/*
				// we were given overrides? If so, lets use them
				if c.Overrides != nil && len(c.Overrides) > 0 {
					for k, v := range c.Overrides {
						b.override[k] = v
					}
				}
				// were we given boosts? Use 'em ...
				if c.Boost != nil && len(c.Boost) > 0 {
					for k, v := range c.Boost {
						b.boost[k] = v
					}
				}
			*/
		} else {
			b.debugMsg("ERROR", unmarshalErr.Error())
		}
	} else {
		b.debugMsg("ERROR", configErr.Error())
	}
}
