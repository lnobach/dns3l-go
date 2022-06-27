package service

import (
	"fmt"
	"io/ioutil"

	"github.com/dta4/dns3l-go/ca"
	"github.com/dta4/dns3l-go/dns"
	"github.com/dta4/dns3l-go/state"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DNS       *dns.Config                 `yaml:"dns"`
	CA        *ca.Config                  `yaml:"ca"`
	RootZones dns.RootZones               `yaml:"rtzn"`
	DB        *state.SQLDBProviderDefault `yaml:"db"` //SQL hard-coded (only here)

	//RootZoneAllowedCA map[string] //maybe we need this later...
	//CAAllowedRootZones map[string][]dns.RootZone
}

func (c *Config) FromFile(file string) error {
	filebytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return c.FromYamlBytes(filebytes)
}

func (c *Config) FromYamlBytes(bytes []byte) error {
	return yaml.Unmarshal(bytes, c)
}

//Must be executed after config struct initialization
func (c *Config) Initialize() error {

	for _, rtzn := range c.RootZones {
		if foo := rtzn.CAs[0]; foo == "*" {
			//all CAs can handle this root zone
			for _, ca := range c.CA.Providers {
				ca.AddAllowedRootZone(rtzn)
			}
			continue
		}
		for _, caID := range rtzn.CAs {
			ca, exists := c.CA.Providers[caID]
			if !exists {
				return fmt.Errorf("CA '%s' has not been configured", caID)
			}
			ca.AddAllowedRootZone(rtzn)
		}
	}

	err := c.CA.Init(&CAConfigurationContextImpl{
		StateProvider: c.DB,
	})
	if err != nil {
		return err
	}

	return nil

}
