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
	c := &Config{}
	configContents, configErr := config.Home("block", ".yml")
	if configErr == nil {
		unmarshalErr := yaml.Unmarshal([]byte(configContents), &c)
		if unmarshalErr == nil {
			// we were given overrides? If so, lets use them
			if len(c.Overrides) > 0 {
				c.Overrides = make(map[string]string)
				for k, v := range c.Overrides {
					b.override[k] = v
				}
			}
			// were we given boosts? Use 'em ...
			if len(c.Boost) > 0 {
				c.Boost = make(map[string]int)
				for k, v := range c.Boost {
					b.boost[k] = v
				}
			}
		} else {
			b.debugMsg("ERROR", unmarshalErr.Error())
		}
	} else {
		b.debugMsg("ERROR", configErr.Error())
	}
}
