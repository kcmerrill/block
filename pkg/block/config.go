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
			b.boost = c.Boost
			b.override = c.Overrides
		} else {
			b.debugMsg("ERROR", unmarshalErr.Error())
		}
	} else {
		b.debugMsg("ERROR", configErr.Error())
	}
}
