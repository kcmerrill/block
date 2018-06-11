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
	c := Config{}
	configContents, _ := config.Home("block", ".yml")
	unmarshalErr := yaml.Unmarshal(configContents, &c)
	if unmarshalErr == nil {
		// we were given overrides? If so, lets use them
		for k, v := range c.Overrides {
			b.Overrides[k] = v
		}
		// were we given boosts? Use 'em ...
		for k, v := range c.Boost {
			b.Boosted[k] = v
		}
	}
}
