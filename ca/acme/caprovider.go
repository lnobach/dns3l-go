package acme

import (
	"errors"

	castate "github.com/dta4/dns3l-go/ca/state"
	"github.com/dta4/dns3l-go/ca/types"
)

type CAProvider struct {
	engine *Engine
	c      *Config
}

func (p *CAProvider) GetInfo() *types.CAProviderInfo {

	return &types.CAProviderInfo{
		Name:        p.c.Name,
		Type:        p.c.CAType,
		Description: "foo <TODO where does this come from?>",
		LogoPath:    "bar <TODO where does this come from?>",
		URL:         p.c.URL,
		Roots:       p.c.Roots,
		IsAcme:      true,
	}

}

func (p *CAProvider) Init(c types.ProviderConfigurationContext) error {

	smgr, err := makeACMEStateManager(c)

	if err != nil {
		return err
	}

	p.engine = &Engine{
		CAID:    c.GetCAID(),
		Conf:    p.c,
		Context: c,
		State:   smgr,
	}

	log.Debugf("ACME CA provider initialized.")

	return nil

}

func makeACMEStateManager(c types.ProviderConfigurationContext) (ACMEStateManager, error) {
	switch sprovinst := c.GetStateMgr().(type) {
	case *castate.CAStateManagerSQL:
		return &ACMEStateManagerSQL{c.GetCAID(), sprovinst.Prov}, nil
	default:
		return nil, errors.New("only supporting SQL DB providers at the moment")
	}
}

func (p *CAProvider) AddAllowedRootZone() int {

	return 42 //TODO

}

func (p *CAProvider) GetTotalValid() int {

	return 42 //TODO

}

func (p *CAProvider) GetTotalIssued() int {

	return 68 //TODO

}

func (p *CAProvider) IsEnabled() bool {

	return true //TODO

}

func (p *CAProvider) ClaimCertificate(cinfo *types.CertificateClaimInfo) error {

	acmeuser := "acme-" + cinfo.Name

	return p.engine.TriggerUpdate(acmeuser, cinfo.Name, cinfo.SubjectAltNames,
		"leo@nobach.net", //TODO: clarify where to get this from.
		"anonymous")      //TODO: fix when auth is ready.

}

func (p *CAProvider) CleanupBeforeDeletion(keyID string) error {

	//TODO maybe clean up orphaned ACME keys here
	return nil

}
